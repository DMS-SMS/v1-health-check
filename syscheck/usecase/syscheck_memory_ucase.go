// Create file in v.1.0.0
// syscheck_memory_ucase.go is file that define usecase implementation about syscheck memory domain
// memory check usecase struct embed systemCheckUsecaseComponent struct in ./syscheck.go file

package usecase

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/google/uuid"
	"github.com/inhies/go-bytesize"
	"github.com/pkg/errors"
	"sync"

	"github.com/DMS-SMS/v1-health-check/domain"
)

// memoryCheckStatus is type to int constant represent current memory check process status
type memoryCheckStatus int
const (
	memoryStatusHealthy    memoryCheckStatus = iota // represent memory check status is healthy
	memoryStatusWarning                             // represent memory check status is warning now
	memoryStatusRecovering                          // represent it's recovering memory status now
	memoryStatusUnhealthy                           // represent memory check status is unhealthy
)

// memoryCheckUsecase implement MemoryCheckUsecase interface in domain and used in delivery layer
type memoryCheckUsecase struct {
	// myCfg is used for getting memory check usecase config
	myCfg memoryCheckUsecaseConfig

	// historyRepo is used for store memory check history and injected from outside
	historyRepo domain.MemoryCheckHistoryRepository

	// slackChat is used for agent slack API about chatting
	slackChatAgency slackChatAgency

	// memorySysAgency is used as agency about memory system command
	memorySysAgency memorySysAgency

	// dockerAgency is used as agency about docker command
	dockerAgency dockerAgency

	// status represent current process status of memory health check
	status memoryCheckStatus

	// mutex help to prevent race condition when set status field value
	mutex sync.Mutex
}

// memoryCheckUsecaseConfig is the config getter interface for memory check usecase
type memoryCheckUsecaseConfig interface {
	// get common config method from embedding systemCheckUsecaseComponentConfig
	systemCheckUsecaseComponentConfig

	// MemoryWarningUsage method returns bytesize.ByteSize represent memory warning usage
	MemoryWarningUsage() bytesize.ByteSize

	// MemoryMaximumUsage method returns bytesize.ByteSize represent memory maximum usage
	MemoryMaximumUsage() bytesize.ByteSize

	// MemoryMinimumUsageToRemove method returns bytesize.ByteSize represent memory minimum usage to remove
	MemoryMinimumUsageToRemove() bytesize.ByteSize
}

// memorySysAgency is agency that agent various command about memory system
type memorySysAgency interface {
	// GetTotalSystemMemoryUsage return total memory usage as bytesize in system
	GetTotalSystemMemoryUsage() (usage bytesize.ByteSize, err error)

	// CalculateContainersMemoryUsage calculate container memory usage & return result interface implementation
	CalculateContainersMemoryUsage() (result interface {
		// TotalMemoryUsage return total memory usage in docker containers
		TotalMemoryUsage() (usage bytesize.ByteSize)

		// MostConsumerExceptFor return container consume the most memory except container names received from param
		MostConsumerExceptFor(names []string) (id, name string, usage bytesize.ByteSize)
	}, err error)
}

// NewMemoryCheckUsecase function return memoryCheckUsecase ptr instance after initializing
func NewMemoryCheckUsecase(
	cfg memoryCheckUsecaseConfig,
	mhr domain.MemoryCheckHistoryRepository,
	sca slackChatAgency,
	msa memorySysAgency,
	da dockerAgency,
) domain.MemoryCheckUseCase {
	return &memoryCheckUsecase{
		// initialize field with parameter received from caller
		myCfg:           cfg,
		historyRepo:     mhr,
		slackChatAgency: sca,
		memorySysAgency: msa,
		dockerAgency:    da,

		// initialize field with default value
		status: memoryStatusHealthy,
		mutex:  sync.Mutex{},
	}
}

// CheckMemory check memory health with CheckMemory method & store check history in repository
// Implement CheckMemory method of domain.MemoryCheckUseCase interface
func (mu *memoryCheckUsecase) CheckMemory(ctx context.Context) error {
	history := mu.checkMemory(ctx)

	if b, err := mu.historyRepo.Store(history); err != nil {
		return errors.Wrapf(err, "failed to store memory check history, response: %s", string(b))
	}

	return nil
}

// method with below logic about handling health check process according to current memory check status
// 0 : 정상적으로 인지된 상태 (상태 확인 수행)
// 0 -> 1 : 메모리 사용량이 Warning 수치보다 높아짐 (경고 상태 알림 발행)
// 1 -> 0 : 메모리 사용량 정상 수치로 복귀 (경고 상태 해제 알림 발행)
// (0 or 1) -> 2 : 메모리 기준 프로비저닝 실행 (상태 회복중 상테 알림 발행)
// 2 : 메모리 프로비저닝 실행중 (상태 확인 수행 X)
// 2 -> 0 : 메모리 프로비저닝으로 인해 상태 회복 완료 (상태 회복 성공 알림 발행)
// 2 -> 3 : 메모리 프로비저닝을 해도 상태 회복 X (상태 회복 불가능 상태 알림 발행)
// 3 : 관리자가 직접 확인해야함 (상태 확인 수행 X)
// 3 -> 0 : 관리자 직접 상태 회복 완료 (상태 회복 알림 발행)
func (mu *memoryCheckUsecase) checkMemory(ctx context.Context) (history *domain.MemoryCheckHistory) {
	_uuid := uuid.New().String()
	history = new(domain.MemoryCheckHistory)
	history.FillPrivateComponent()
	history.UUID = _uuid

	_totalUsage, err := mu.memorySysAgency.GetTotalSystemMemoryUsage()
	if err != nil {
		history.ProcessLevel.Set(errorLevel)
		history.SetError(errors.Wrap(err, "failed to get total system memory usage"))
		return
	}
	history.TotalUsageMemory = _totalUsage
	var totalUsage = bytesizeComparator{_totalUsage}

	switch mu.status {
	case memoryStatusHealthy:
		break
	case memoryStatusWarning:
		if totalUsage.isLessThan(mu.myCfg.MemoryWarningUsage()) {
			mu.setStatus(memoryStatusHealthy)
		}
	case memoryStatusRecovering:
		history.ProcessLevel.Set(recoveringLevel)
		history.Message = "provisioning memory is already on process using docker"
		return
	case memoryStatusUnhealthy:
		if totalUsage.isLessThan(mu.myCfg.MemoryMaximumUsage()) {
			mu.setStatus(memoryStatusHealthy)
			history.ProcessLevel.Set(recoveredLevel)
			history.Message = "memory check is recovered to be healthy"
			msg := fmt.Sprintf("!memory check recovered to health! current memory usage - %s", totalUsage)
			_, _, _ = mu.slackChatAgency.SendMessage("heart", msg, _uuid)
		} else {
			history.ProcessLevel.Set(unhealthyLevel)
			history.Message = "memory check is unhealthy now"
		}
		return
	}

	if totalUsage.isMoreThan(mu.myCfg.MemoryMaximumUsage()) {
		mu.setStatus(memoryStatusRecovering)
		history.ProcessLevel.Set(weakDetectedLevel)
		msg := fmt.Sprintf("!memory check weak detected! start to provision memory (current memory usage - %s)", totalUsage)
		history.SetAlarmResult(mu.slackChatAgency.SendMessage("pill", msg, _uuid))

		result, err := mu.memorySysAgency.CalculateContainersMemoryUsage()
		if err != nil {
			mu.setStatus(memoryStatusUnhealthy)
			history.ProcessLevel.Append(errorLevel)
			msg := "!memory check error occurred! failed to calculate container memory, please check for yourself"
			_, _, _ = mu.slackChatAgency.SendMessage("anger", msg, _uuid)
			history.SetError(errors.Wrap(err, "failed to calculate containers memory usage"))
			return
		}
		history.DockerUsageMemory = result.TotalMemoryUsage()

		id, name, _usage := result.MostConsumerExceptFor(requiredContainers)
		history.MostMemoryConsumeContainer = name
		usage := bytesizeComparator{_usage}

		if usage.isLessThan(mu.myCfg.MemoryMinimumUsageToRemove()) {
			mu.setStatus(memoryStatusUnhealthy)
			msg := "!memory check error occurred! memory usage is too small to remove, please check for yourself"
			_, _, _ = mu.slackChatAgency.SendMessage("anger", msg, _uuid)
			history.SetError(errors.New("memory usage is too small to remove"))
			return
		}

		if err := mu.dockerAgency.RemoveContainer(id, types.ContainerRemoveOptions{Force: true}); err != nil {
			mu.setStatus(memoryStatusUnhealthy)
			history.ProcessLevel.Append(errorLevel)
			msg := "!memory check error occurred! failed to remove container, please check for yourself"
			_, _, _ = mu.slackChatAgency.SendMessage("anger", msg, _uuid)
			history.SetError(errors.Wrap(err, "failed to remove container"))
			return
		} else {
			history.TemporaryFreeMemory = usage.ByteSize
			history.Message = "removed most memory consumed container as memory usage is over than maximum"
		}

		_againTotalUsage, err := mu.memorySysAgency.GetTotalSystemMemoryUsage()
		if err != nil {
			mu.setStatus(memoryStatusUnhealthy)
			history.ProcessLevel.Append(errorLevel)
			msg := "!memory check error occurred! failed to again calculate container memory, please check for yourself"
			_, _, _ = mu.slackChatAgency.SendMessage("broken_heart", msg, _uuid)
			history.SetError(errors.Wrap(err, "failed to again calculate containers memory usage"))
			return
		}
		var againTotalUsage = bytesizeComparator{_againTotalUsage}

		if againTotalUsage.isLessThan(mu.myCfg.MemoryMaximumUsage()) {
			mu.setStatus(memoryStatusHealthy)
			msg := fmt.Sprintf("!memory check is healthy! current memory usage - %s", againTotalUsage)
			_, _, _ = mu.slackChatAgency.SendMessage("heart", msg, _uuid)
		} else {
			mu.setStatus(memoryStatusUnhealthy)
			msg := "!memory check has deteriorated! please check for yourself"
			_, _, _ = mu.slackChatAgency.SendMessage("broken_heart", msg, _uuid)
		}
	} else if totalUsage.isMoreThan(mu.myCfg.MemoryWarningUsage()) {
		history.ProcessLevel.Set(warningLevel)
		history.Message = "memory check is warning now, but not weak yet"
		if mu.status != memoryStatusWarning {
			mu.setStatus(memoryStatusWarning)
			msg := fmt.Sprintf("!memory check warning! current memory usage - %s", totalUsage)
			history.SetAlarmResult(mu.slackChatAgency.SendMessage("warning", msg, _uuid))
		}
	} else {
		history.ProcessLevel.Set(healthyLevel)
		history.Message = "memory system is healthy now"
	}

	return
}

// setStatus set status field value using mutex Lock & Unlock
func (mu *memoryCheckUsecase) setStatus(status memoryCheckStatus) {
	mu.mutex.Lock()
	defer mu.mutex.Unlock()
	mu.status = status
}
