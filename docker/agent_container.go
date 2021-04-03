// Create file in v.1.0.0
// agent_container.go is file that define method of dockerAgent that agent command about docker container
// For example in disk command, there are remove or restart container, etc ...

package docker

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/inhies/go-bytesize"
	"github.com/pkg/errors"
	"strings"
)

// GetContainerWithServiceName return container which is instance of received service name
func (da *dockerAgent) GetContainerWithServiceName(srv string) (interface {
	ID() string                     // get id of container
	MemoryUsage() bytesize.ByteSize // get memory usage of container
}, error) {
	var (
		ctx = context.Background()
	)
	
	containers, err := da.dkrCli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get container list from docker")
	}

	for _, ctn := range containers {
		// ex) [/DSM_SMS_api-gateway.1.mod9z6n0hey4n6topphc2700r] -> [DSM_SMS_api-gateway 1 mod9z6n0hey4n6topphc2700r]
		sep := strings.Split(strings.TrimPrefix(ctn.Names[0], "/"), ".")
		if len(sep) == 0 || sep[0] != srv {
			continue
		}

		stats, err := da.dkrCli.ContainerStats(ctx, ctn.ID, false)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get container stats from docker")
		}

		v := &types.StatsJSON{}
		if err = json.NewDecoder(stats.Body).Decode(v); err != nil {
			return nil, errors.Wrap(err, "failed to decode stats response body to struct")
		}

		size, err := getMemoryUsageSizeFrom(v)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get memory usage size from Stats")
		}

		return container{
			id:          ctn.ID,
			memoryUsage: size,
		}, nil
	}
	
	return nil, errors.New("container with that service name does't exist")
}

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
