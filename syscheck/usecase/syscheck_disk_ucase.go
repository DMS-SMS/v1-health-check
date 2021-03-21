// Create file in v.1.0.0
// syscheck_disk_ucase.go is file that define usecase implementation about disk check domain
// disk check usecase struct embed systemCheckUsecaseComponent struct in ./syscheck.go file

package usecase

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/inhies/go-bytesize"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
	"os"
	"sync"

	"github.com/DMS-SMS/v1-health-check/domain"
)

// diskCheckStatus is type to int constant represent current disk check process status
type diskCheckStatus int
const (
	healthyStatus diskCheckStatus = iota // represent disk check status is healthy
	recoveringStatus                     // represent it's recovering disk status now
	unhealthyStatus                      // represent disk check status is unhealthy
)

// diskCheckUsecase implement DiskCheckUsecase interface in domain and used in delivery layer
type diskCheckUsecase struct {
	// myCfg is used for getting disk check usecase config
	myCfg diskCheckUsecaseConfig

	// historyRepo is used for store disk check history and injected from outside
	historyRepo domain.DiskCheckHistoryRepository

	// dockerCli is docker client to call docker agent API
	dockerCli *client.Client

	// slackChat is used for agent slack API about chatting
	slackChatAgency slackChatAgency

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

// NewDiskCheckUsecase function return diskCheckUsecase ptr instance with initializing
func NewDiskCheckUsecase(cfg diskCheckUsecaseConfig, hr domain.DiskCheckHistoryRepository, dc *client.Client, sca slackChatAgency) *diskCheckUsecase {
	return &diskCheckUsecase{
		// initialize field with parameter received from caller
		myCfg:           cfg,
		historyRepo:     hr,
		dockerCli:       dc,
		slackChatAgency: sca,

		// initialize field with default value
		status: healthyStatus,
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

	_cap, err := getRemainDiskCapacity()
	if err != nil {
		err = errors.Wrap(err, "failed to get disk capacity")
		history.ProcessLevel = errorLevel.String()
		history.SetError(err)
		return
	}
	history.DiskCapacity = _cap

	switch du.status {
	case healthyStatus:
		break
	case recoveringStatus:
		history.ProcessLevel = recoveringLevel.String()
		history.Message = "pruning docker system is already on process"
		return
	case unhealthyStatus:
		if du.isMinCapacityLessThan(_cap) {
			du.setStatus(healthyStatus)
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
		du.setStatus(recoveringStatus)
		history.ProcessLevel = weakDetectedLevel.String()
		msg := "!weak detected in disk check! start to prune docker system"
		history.SetAlarmResult(du.slackChatAgency.SendMessage("warning", msg, _uuid))

		if r, err := du.pruneDockerSystem(); err != nil {
			msg := "!disk check error occurred! failed to prune docker system"
			_, _, _ = du.slackChatAgency.SendMessage("anger", msg, _uuid)
			err = errors.Wrap(err, "failed to prune docker system")
			history.SetError(err)
		} else {
			history.ReclaimedSize = r
			history.Message = "pruned docker system as current disk capacity is less than the minimum"
		}

		if _cap, _ = getRemainDiskCapacity(); du.isMinCapacityLessThan(_cap) {
			du.setStatus(healthyStatus)
			msg := fmt.Sprintf("!disk check is healthy by pruning! remain capacity - %s", _cap.String())
			_, _, _ = du.slackChatAgency.SendMessage("heart", msg, _uuid)
		} else {
			du.setStatus(unhealthyStatus)
			msg := "!disk check has deteriorated! please check for yourself"
			_, _, _ = du.slackChatAgency.SendMessage("broken_heart", msg, _uuid)
		}
	} else {
		history.ProcessLevel = healthyLevel.String()
		history.Message = "disk system is healthy now"
	}

	return
}

// pruneDockerSystem prune docker system(build cache, containers, images, networks) and return reclaimed space size
func (du *diskCheckUsecase) pruneDockerSystem() (reclaimed bytesize.ByteSize, err error) {
	var (
		ctx = context.Background()
		args = filters.Args{}
	)

	if report, pruneErr := du.dockerCli.BuildCachePrune(ctx, types.BuildCachePruneOptions{}); pruneErr != nil {
		err = errors.Wrap(pruneErr, "failed to prune build cache in docker")
		return
	} else {
		reclaimed = bytesize.ByteSize(uint64(reclaimed) + report.SpaceReclaimed)
	}

	if report, pruneErr := du.dockerCli.ContainersPrune(ctx, args); pruneErr != nil {
		err = errors.Wrap(pruneErr, "failed to prune containers in docker")
		return
	} else {
		reclaimed = bytesize.ByteSize(uint64(reclaimed) + report.SpaceReclaimed)
	}

	if report, pruneErr := du.dockerCli.ImagesPrune(ctx, args); pruneErr != nil {
		err = errors.Wrap(pruneErr, "failed to prune image in docker")
		return
	} else {
		reclaimed = bytesize.ByteSize(uint64(reclaimed) + report.SpaceReclaimed)
	}

	if _, pruneErr := du.dockerCli.NetworksPrune(ctx, args); pruneErr != nil {
		err = errors.Wrap(pruneErr, "failed to prune networks in docker")
		return
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

// getRemainDiskCapacity returns remain disk capacity as bytesize.Bytesize
func getRemainDiskCapacity() (size bytesize.ByteSize, err error) {
	var stat unix.Statfs_t

	wd, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "failed to call os.Getwd")
		return
	}

	if err = unix.Statfs(wd, &stat); err != nil {
		err = errors.Wrap(err, "failed to call unix.Statfs")
		return
	}

	// Available blocks * size per block = available space in bytes
	size = bytesize.New(float64(stat.Bavail * uint64(stat.Bsize)))
	return
}
