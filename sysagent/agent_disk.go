// Create file in v.1.0.0
// agent_disk.go is file that define method of sysAgent that agent command about disk
// For example in disk command, there are get remaining capacity, prune disk, etc ...

package sysagent

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/inhies/go-bytesize"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
	"os"
)

// GetRemainDiskCapacity return remain disk capacity expressed in bytesize package
func (sa *sysAgent) GetRemainDiskCapacity() (size bytesize.ByteSize, err error) {
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

// PruneDockerSystem prune docker system(build cache, containers, images, networks) and return reclaimed space size
func (sa *sysAgent) PruneDockerSystem() (reclaimed bytesize.ByteSize, err error) {
	var (
		ctx = context.Background()
		args = filters.Args{}
	)

	if report, pruneErr := sa.dockerCli.BuildCachePrune(ctx, types.BuildCachePruneOptions{}); pruneErr != nil {
		err = errors.Wrap(pruneErr, "failed to prune build cache in docker")
		return
	} else {
		reclaimed = bytesize.ByteSize(uint64(reclaimed) + report.SpaceReclaimed)
	}

	if report, pruneErr := sa.dockerCli.ContainersPrune(ctx, args); pruneErr != nil {
		err = errors.Wrap(pruneErr, "failed to prune containers in docker")
		return
	} else {
		reclaimed = bytesize.ByteSize(uint64(reclaimed) + report.SpaceReclaimed)
	}

	if report, pruneErr := sa.dockerCli.ImagesPrune(ctx, args); pruneErr != nil {
		err = errors.Wrap(pruneErr, "failed to prune image in docker")
		return
	} else {
		reclaimed = bytesize.ByteSize(uint64(reclaimed) + report.SpaceReclaimed)
	}

	if _, pruneErr := sa.dockerCli.NetworksPrune(ctx, args); pruneErr != nil {
		err = errors.Wrap(pruneErr, "failed to prune networks in docker")
		return
	}

	return
}
