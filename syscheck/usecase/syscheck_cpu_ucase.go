// Create file in v.1.0.0
// syscheck_cpu_ucase.go is file that define usecase implementation about syscheck cpu domain
// cpu check usecase struct embed systemCheckUsecaseComponent struct in ./syscheck.go file

package usecase

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"sync"

	"github.com/DMS-SMS/v1-health-check/domain"
)

// cpuCheckStatus is type to int constant represent current cpu check process status
type cpuCheckStatus int

const (
	cpuStatusHealthy    cpuCheckStatus = iota // represent cpu check status is healthy
	cpuStatusWarning                          // represent cpu check status is warning now
	cpuStatusRecovering                       // represent it's recovering cpu status now
	cpuStatusUnhealthy                        // represent cpu check status is unhealthy
)

// cpuCheckUsecase implement CPUCheckUsecase interface in domain and used in delivery layer
type cpuCheckUsecase struct {
	// myCfg is used for getting cpu check usecase config
	myCfg cpuCheckUsecaseConfig

	// historyRepo is used for store cpu check history and injected from outside
	historyRepo domain.CPUCheckHistoryRepository

	// slackChat is used for agent slack API about chatting
	slackChatAgency slackChatAgency

	// cpuSysAgency is used as agency about cpu system command
	cpuSysAgency cpuSysAgency

	// dockerAgency is used as agency about docker command
	dockerAgency dockerAgency

	// status represent current process status of cpu health check
	status cpuCheckStatus

	// mutex help to prevent race condition when set status field value
	mutex sync.Mutex
}

// cpuCheckUsecaseConfig is the config getter interface for cpu check usecase
type cpuCheckUsecaseConfig interface {
	// get common config method from embedding systemCheckUsecaseComponentConfig
	systemCheckUsecaseComponentConfig

	// CPUWarningUsage method returns float64 represent cpu warning usage
	CPUWarningUsage() float64

	// CPUMaximumUsage method returns float64 represent cpu maximum usage
	CPUMaximumUsage() float64

	// CPUMinimumUsageToRemove method returns float64 represent cpu minimum usage to remove
	CPUMinimumUsageToRemove() float64
}

// cpuSysAgency is agency that agent various command about cpu system
type cpuSysAgency interface {
	// GetTotalSystemCPUUsage return total cpu usage as core count in system
	GetTotalSystemCPUUsage() (usage float64, err error)

	// CalculateContainersCPUUsage calculate container cpu usage & return result interface implementation
	CalculateContainersCPUUsage() (result interface {
		// TotalCPUUsage return total cpu usage in docker containers
		TotalCPUUsage() (usage float64)

		// MostConsumerExceptFor return container consume the most CPU except container names received from param
		MostConsumerExceptFor(names []string) (id, name string, usage float64)
	}, err error)
}

// NewCPUCheckUsecase function return cpuCheckUsecase ptr instance after initializing
func NewCPUCheckUsecase(
	cfg cpuCheckUsecaseConfig,
	chr domain.CPUCheckHistoryRepository,
	sca slackChatAgency,
	csa cpuSysAgency,
	da dockerAgency,
) domain.CPUCheckUseCase {
	return &cpuCheckUsecase{
		// initialize field with parameter received from caller
		myCfg:           cfg,
		historyRepo:     chr,
		slackChatAgency: sca,
		cpuSysAgency:    csa,
		dockerAgency:    da,

		// initialize field with default value
		status: cpuStatusHealthy,
		mutex:  sync.Mutex{},
	}
}

// CheckCPU check cpu health with checkCPU method & store check history in repository
// Implement CheckCPU method of domain.CPUCheckUseCase interface
func (cu *cpuCheckUsecase) CheckCPU(ctx context.Context) error {
	history := cu.checkCPU(ctx)

	if b, err := cu.historyRepo.Store(history); err != nil {
		return errors.Wrapf(err, "failed to store cpu check history, response: %s", string(b))
	}

	return nil
}

// method with below logic about handling health check process according to current cpu check status
// 0 : 정상적으로 인지된 상태 (상태 확인 수행)
// 0 -> 1 : CPU 사용량이 Warning 수치보다 높아짐 (경고 상태 알림 발행)
// 1 -> 0 : CPU 사용량 정상 수치로 복귀 (경고 상태 해제 알림 발행)
// (0 or 1) -> 2 : CPU 기준 프로비저닝 실행 (상태 회복중 상테 알림 발행)
// 2 : CPU 프로비저닝 실행중 (상태 확인 수행 X)
// 2 -> 0 : CPU 프로비저닝으로 인해 상태 회복 완료 (상태 회복 성공 알림 발행)
// 2 -> 3 : CPU 프로비저닝을 해도 상태 회복 X (상태 회복 불가능 상태 알림 발행)
// 3 : 관리자가 직접 확인해야함 (상태 확인 수행 X)
// 3 -> 0 : 관리자 직접 상태 회복 완료 (상태 회복 알림 발행)
func (cu *cpuCheckUsecase) checkCPU(ctx context.Context) (history *domain.CPUCheckHistory) {
	_uuid := uuid.New().String()
	history = new(domain.CPUCheckHistory)
	history.FillPrivateComponent()
	history.UUID = _uuid

	_totalUsage, err := cu.cpuSysAgency.GetTotalSystemCPUUsage()
	if err != nil {
		history.ProcessLevel.Set(errorLevel)
		history.SetError(errors.Wrap(err, "failed to get total system cpu usage"))
		msg := "!cpu check error occurred! unable to get total cpu usage"
		history.SetAlarmResult(cu.slackChatAgency.SendMessage("x", msg, _uuid))
		return
	}
	history.TotalUsageCore = _totalUsage
	var totalUsage = float64Comparator{V: _totalUsage}

	switch cu.status {
	case cpuStatusHealthy:
		break
	case cpuStatusWarning:
		if totalUsage.isLessThan(cu.myCfg.CPUWarningUsage()) {
			cu.setStatus(cpuStatusHealthy)
		}
	case cpuStatusRecovering:
		history.ProcessLevel.Set(recoveringLevel)
		history.Message = "provisioning CPU is already on process using docker"
		return
	case cpuStatusUnhealthy:
		if totalUsage.isLessThan(cu.myCfg.CPUMaximumUsage()) {
			cu.setStatus(cpuStatusHealthy)
			history.ProcessLevel.Set(recoveredLevel)
			history.Message = "cpu check is recovered to be healthy"
			msg := fmt.Sprintf("!cpu check recovered to health! current cpu usage - %.02f", totalUsage.V)
			_, _, _ = cu.slackChatAgency.SendMessage("heart", msg, _uuid)
		} else {
			history.ProcessLevel.Set(unhealthyLevel)
			history.Message = "cpu check is unhealthy now"
		}
		return
	}

	if totalUsage.isMoreThan(cu.myCfg.CPUMaximumUsage()) {
		cu.setStatus(cpuStatusRecovering)
		history.ProcessLevel.Set(weakDetectedLevel)
		msg := fmt.Sprintf("!cpu check weak detected! start to provision CPU (current cpu usage - %.02f)", totalUsage.V)
		history.SetAlarmResult(cu.slackChatAgency.SendMessage("pill", msg, _uuid))

		result, err := cu.cpuSysAgency.CalculateContainersCPUUsage()
		if err != nil {
			cu.setStatus(cpuStatusUnhealthy)
			history.ProcessLevel.Append(errorLevel)
			msg := "!cpu check error occurred! failed to calculate container cpu, please check for yourself"
			_, _, _ = cu.slackChatAgency.SendMessage("anger", msg, _uuid)
			history.SetError(errors.Wrap(err, "failed to calculate containers cpu usage"))
			return
		}
		history.DockerUsageCore = result.TotalCPUUsage()

		id, name, _usage := result.MostConsumerExceptFor(requiredContainers)
		history.MostCPUConsumeContainer = name
		var usage = float64Comparator{V: _usage}

		if usage.isLessThan(cu.myCfg.CPUMinimumUsageToRemove()) {
			cu.setStatus(cpuStatusUnhealthy)
			msg := "!cpu check error occurred! cpu usage is too small to remove, please check for yourself"
			_, _, _ = cu.slackChatAgency.SendMessage("anger", msg, _uuid)
			history.SetError(errors.New("cpu usage is too small to remove"))
			return
		}

		if err := cu.dockerAgency.RemoveContainer(id, types.ContainerRemoveOptions{Force: true}); err != nil {
			cu.setStatus(cpuStatusUnhealthy)
			history.ProcessLevel.Append(errorLevel)
			msg := "!cpu check error occurred! failed to remove container, please check for yourself"
			_, _, _ = cu.slackChatAgency.SendMessage("anger", msg, _uuid)
			history.SetError(errors.Wrap(err, "failed to remove container"))
			return
		} else {
			history.TemporaryFreeCore = usage.V
			history.Message = "removed most cpu consumed container as cpu usage is over than maximum"
		}

		_againTotalUsage, err := cu.cpuSysAgency.GetTotalSystemCPUUsage()
		if err != nil {
			cu.setStatus(cpuStatusUnhealthy)
			history.ProcessLevel.Append(errorLevel)
			msg := "!cpu check error occurred! failed to again calculate container cpu, please check for yourself"
			_, _, _ = cu.slackChatAgency.SendMessage("broken_heart", msg, _uuid)
			history.SetError(errors.Wrap(err, "failed to again calculate containers cpu usage"))
			return
		}
		var againTotalUsage = float64Comparator{V: _againTotalUsage}

		if againTotalUsage.isLessThan(cu.myCfg.CPUMaximumUsage()) {
			cu.setStatus(cpuStatusHealthy)
			msg := fmt.Sprintf("!cpu check is healthy! current cpu usage - %.02f", againTotalUsage.V)
			_, _, _ = cu.slackChatAgency.SendMessage("heart", msg, _uuid)
		} else {
			cu.setStatus(cpuStatusUnhealthy)
			msg := "!cpu check has deteriorated! please check for yourself"
			_, _, _ = cu.slackChatAgency.SendMessage("broken_heart", msg, _uuid)
		}
	} else if totalUsage.isMoreThan(cu.myCfg.CPUWarningUsage()) {
		history.ProcessLevel.Set(warningLevel)
		history.Message = "cpu check is warning now, but not weak yet"
		if cu.status != cpuStatusWarning {
			cu.setStatus(cpuStatusWarning)
			msg := fmt.Sprintf("!cpu check warning! current cpu usage - %.02f", totalUsage.V)
			history.SetAlarmResult(cu.slackChatAgency.SendMessage("warning", msg, _uuid))
		}
	} else {
		history.ProcessLevel.Set(healthyLevel)
		history.Message = "cpu system is healthy now"
	}

	return
}

// setStatus set status field value using mutex Lock & Unlock
func (cu *cpuCheckUsecase) setStatus(status cpuCheckStatus) {
	cu.mutex.Lock()
	defer cu.mutex.Unlock()
	cu.status = status
}
