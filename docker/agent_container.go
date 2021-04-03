// Create file in v.1.0.0
// agent_container.go is file that define method of dockerAgent that agent command about docker container
// For example in disk command, there are remove or restart container, etc ...

package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/inhies/go-bytesize"
	"github.com/pkg/errors"
)

// RemoveContainer remove container with id & option (auto created from docker swarm if exists)
func (da *dockerAgent) RemoveContainer(containerID string, options types.ContainerRemoveOptions) error {
	var (
		ctx = context.Background()
	)

	return errors.Wrap(da.dkrCli.ContainerRemove(ctx, containerID, options), "failed to call ContainerRemove")
}

// getMemoryUsageSizeFrom return memory cpu usage as bytesize.Bytesize type from types.StatsJson struct
func getMemoryUsageSizeFrom(v *types.StatsJSON) (size bytesize.ByteSize, err error) {
	size = bytesize.ByteSize(v.MemoryStats.Usage)

	if b, ok := v.MemoryStats.Stats["inactive_anon"]; ok {
		size -= bytesize.ByteSize(b)
	} else {
		err = errors.Wrap(err, "inactive_anon is not exist in MemoryStats.Stats")
		return
	}

	if b, ok := v.MemoryStats.Stats["inactive_file"]; ok {
		size -= bytesize.ByteSize(b)
	} else {
		err = errors.Wrap(err, "inactive_file is not exist in MemoryStats.Stats")
		return
	}

	return
}

// container is struct having inform about container, and implementation of GetContainerWithServiceName return type interface
type container struct {
	id          string
	memoryUsage bytesize.ByteSize
}

// define return field value methods in container
func (c container) ID() string                     { return c.id }
func (c container) MemoryUsage() bytesize.ByteSize { return c.memoryUsage }
