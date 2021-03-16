// Create file in v.1.0.0
// syscheck_disk_ucase.go is file that define usecase implementation about disk check domain
// disk check usecase struct embed systemCheckUsecaseComponent struct in ./syscheck.go file

package usecase

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/inhies/go-bytesize"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
	"os"

	"github.com/DMS-SMS/v1-health-check/domain"
)

// diskCheckUsecase implement DiskCheckUsecase interface in domain and used in delivery layer
type diskCheckUsecase struct {
	// myCfg is used for getting disk check usecase config
	myCfg diskCheckUsecaseConfig

	// historyRepo is used for store disk check history and injected from outside
	historyRepo domain.DiskCheckHistoryRepository

	// dkrCli is docker client to call docker agent API
	dkrCli *client.Client
}

// diskCheckUsecaseConfig is the config getter interface for disk check usecase
type diskCheckUsecaseConfig interface {
	// get common config method from embedding systemCheckUsecaseComponentConfig
	systemCheckUsecaseComponentConfig

	// DiskMinCapacity method returns byte size represent disk minimum capacity
	DiskMinCapacity() bytesize.ByteSize
}

// NewDiskCheckUsecase function return diskCheckUsecase ptr instance with initializing
func NewDiskCheckUsecase(cfg diskCheckUsecaseConfig, hr domain.DiskCheckHistoryRepository, cli *client.Client) *diskCheckUsecase {
	return &diskCheckUsecase{
		myCfg:       cfg,
		historyRepo: hr,
		dkrCli:      cli,
	}
}

// pruneDockerSystem prune docker system(build cache, containers, images, networks) and return reclaimed space size
func (du *diskCheckUsecase) pruneDockerSystem() (reclaimed bytesize.ByteSize, err error) {
	var (
		ctx = context.Background()
		args = filters.Args{}
	)

	if report, pruneErr := du.dkrCli.BuildCachePrune(ctx, types.BuildCachePruneOptions{}); pruneErr != nil {
		err = errors.Wrap(pruneErr, "failed to prune build cache in docker")
		return
	} else {
		reclaimed = bytesize.ByteSize(uint64(reclaimed) + report.SpaceReclaimed)
	}

	if report, pruneErr := du.dkrCli.ContainersPrune(ctx, args); pruneErr != nil {
		err = errors.Wrap(pruneErr, "failed to prune containers in docker")
		return
	} else {
		reclaimed = bytesize.ByteSize(uint64(reclaimed) + report.SpaceReclaimed)
	}

	if report, pruneErr := du.dkrCli.ImagesPrune(ctx, args); pruneErr != nil {
		err = errors.Wrap(pruneErr, "failed to prune image in docker")
		return
	} else {
		reclaimed = bytesize.ByteSize(uint64(reclaimed) + report.SpaceReclaimed)
	}

	if _, pruneErr := du.dkrCli.NetworksPrune(ctx, args); pruneErr != nil {
		err = errors.Wrap(pruneErr, "failed to prune networks in docker")
		return
	}

	return
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
