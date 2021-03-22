// Create file in v.1.0.0
// syscheck_disk_ucase.go is file that define usecase implementation about disk check domain
// disk check usecase struct embed systemCheckUsecaseComponent struct in ./syscheck.go file

package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/inhies/go-bytesize"
	"github.com/pkg/errors"
	"sync"

	"github.com/DMS-SMS/v1-health-check/domain"
)

// diskCheckStatus is type to int constant represent current disk check process status
type diskCheckStatus int
const (
	diskStatusHealthy    diskCheckStatus = iota // represent disk check status is healthy
	diskStatusRecovering                        // represent it's recovering disk status now
	diskStatusUnhealthy                         // represent disk check status is unhealthy
)

// diskCheckUsecase implement DiskCheckUsecase interface in domain and used in delivery layer
type diskCheckUsecase struct {
	// myCfg is used for getting disk check usecase config
	myCfg diskCheckUsecaseConfig

	// historyRepo is used for store disk check history and injected from outside
	historyRepo domain.DiskCheckHistoryRepository

	// slackChat is used for agent slack API about chatting
	slackChatAgency slackChatAgency

	// diskMeasurer is used for measuring various value about disk
	diskSysAgency diskSysAgency

	// status represent current process status of disk health check
	status diskCheckStatus

	// mutex help to prevent race condition when set status field value
	mutex sync.Mutex
}

// diskCheckUsecaseConfig is the config getter interface for disk check usecase
type diskCheckUsecaseConfig interface {
	// get common config method from embedding systemCheckUsecaseComponentConfig
	systemCheckUsecaseComponentConfig

	// DiskMinCapacity method returns byte size represent disk minimum capacity
	DiskMinCapacity() bytesize.ByteSize
}

// diskSysAgency is agency that agent various command about disk system
type diskSysAgency interface {
	// GetRemainDiskCapacity return remain disk capacity expressed in bytesize package
	GetRemainDiskCapacity() (size bytesize.ByteSize, err error)

	// PruneDockerSystem prune all about docker system and return reclaimed size
	PruneDockerSystem() (reclaimed bytesize.ByteSize, err error)
}

// NewDiskCheckUsecase function return diskCheckUsecase ptr instance with initializing
func NewDiskCheckUsecase(
	cfg diskCheckUsecaseConfig,
	dhr domain.DiskCheckHistoryRepository,
	sca slackChatAgency,
	dsa diskSysAgency,
) domain.DiskCheckUseCase {
	return &diskCheckUsecase{
		// initialize field with parameter received from caller
		myCfg:           cfg,
		historyRepo:     dhr,
		slackChatAgency: sca,
		diskSysAgency:   dsa,

		// initialize field with default value
		status: diskStatusHealthy,
		mutex:  sync.Mutex{},
	}
}

// CheckDisk check disk health with checkDisk method & store check log in repository
// Implement CheckDisk method of domain.DiskCheckUseCase interface
func (du *diskCheckUsecase) CheckDisk(ctx context.Context) error {
	history := du.checkDisk(ctx)

	if b, err := du.historyRepo.Store(history); err != nil {
		return errors.Wrapf(err, "failed to store disk check history, response: %s", string(b))
	}

	return nil
}

// method with below logic about handling health check process according to current disk check status
// 0 : 정상적으로 인지된 상태 (상태 확인 수행)
// 0 -> 1 : Docker Prune 실행 (Docker Prune 알림 발행)
// 1 : Docker Prune 실행중 (상태 확인 수행 X)
// 1 -> 0 : Docker Prune 으로 인해 상태 회복 완료 (상태 회복 알림 발행)
// 1 -> 2 : Docker Prune 을 해도 상태 회복 X (상태 회복 불가능 상태 알림 발행)
// 2 : 관리자가 직접 확인해야함 (상태 확인 수행 X)
// 2 -> 0 : 관리자 직접 상태 회복 완료 (상태 회복 알림 발행)
func (du *diskCheckUsecase) checkDisk(ctx context.Context) (history *domain.DiskCheckHistory) {
	_uuid := uuid.New().String()
	history = new(domain.DiskCheckHistory)
	history.FillPrivateComponent()
	history.UUID = _uuid

	_cap, err := du.diskSysAgency.GetRemainDiskCapacity()
	if err != nil {
		err = errors.Wrap(err, "failed to get disk capacity")
		history.ProcessLevel = errorLevel.String()
		history.SetError(err)
		return
	}
	history.RemainingCap = _cap

	switch du.status {
	case diskStatusHealthy:
		break
	case diskStatusRecovering:
		history.ProcessLevel = recoveringLevel.String()
		history.Message = "pruning docker system is already on process"
		return
	case diskStatusUnhealthy:
		if du.isMinCapacityLessThan(_cap) {
			du.setStatus(diskStatusHealthy)
			history.ProcessLevel = recoveredLevel.String()
			history.Message = "disk check is recovered to be healthy"
			msg := fmt.Sprintf("!disk check recovered to health! remain capacity - %s", _cap.String())
			_, _, _ = du.slackChatAgency.SendMessage("heart", msg, _uuid)
		} else {
			history.ProcessLevel = unhealthyLevel.String()
			history.Message = "disk check is unhealthy now"
		}
		return
	}

	if !du.isMinCapacityLessThan(_cap) {
		du.setStatus(diskStatusRecovering)
		history.ProcessLevel = weakDetectedLevel.String()
		msg := "!weak detected in disk check! start to prune docker system"
		history.SetAlarmResult(du.slackChatAgency.SendMessage("warning", msg, _uuid))

		if r, err := du.diskSysAgency.PruneDockerSystem(); err != nil {
			msg := "!disk check error occurred! failed to prune docker system"
			_, _, _ = du.slackChatAgency.SendMessage("anger", msg, _uuid)
			err = errors.Wrap(err, "failed to prune docker system")
			history.SetError(err)
		} else {
			history.ReclaimedCap = r
			history.Message = "pruned docker system as current disk capacity is less than the minimum"
		}

		if _cap, _ = du.diskSysAgency.GetRemainDiskCapacity(); du.isMinCapacityLessThan(_cap) {
			du.setStatus(diskStatusHealthy)
			msg := fmt.Sprintf("!disk check is healthy by pruning! remain capacity - %s", _cap.String())
			_, _, _ = du.slackChatAgency.SendMessage("heart", msg, _uuid)
		} else {
			du.setStatus(diskStatusUnhealthy)
			msg := "!disk check has deteriorated! please check for yourself"
			_, _, _ = du.slackChatAgency.SendMessage("broken_heart", msg, _uuid)
		}
	} else {
		history.ProcessLevel = healthyLevel.String()
		history.Message = "disk system is healthy now"
	}

	return
}

// isMinCapacityLessThan return bool if disk min capacity is less than parameter
func (du *diskCheckUsecase) isMinCapacityLessThan(_cap bytesize.ByteSize) bool {
	return du.myCfg.DiskMinCapacity() < _cap
}

// setStatus set status field value using mutex Lock & Unlock
func (du *diskCheckUsecase) setStatus(status diskCheckStatus) {
	du.mutex.Lock()
	defer du.mutex.Unlock()
	du.status = status
}
