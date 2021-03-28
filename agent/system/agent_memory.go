// Create file in v.1.0.0
// agent_memory.go is file that define method of sysAgent that agent command about memory
// For example in memory command, there are get total memory usage, calculate container memory usage, etc ...

package system

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/inhies/go-bytesize"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/pkg/errors"
)

// GetTotalSystemMemoryUsage return total memory usage as bytesize in system
func (sa *sysAgent) GetTotalSystemMemoryUsage() (usage bytesize.ByteSize, err error) {
	stats, err := memory.Get()
	if err != nil {
		err = errors.Wrap(err, "failed to get memory stats")
		return
	}

	usage = bytesize.ByteSize(stats.Used)
	return
}

// CalculateContainersCPUUsage calculate memory usage & return calculateContainersMemoryUsageResult
func (sa *sysAgent) CalculateContainersMemoryUsage() (interface {
	TotalMemoryUsage() (usage bytesize.ByteSize)
	MostConsumerExceptFor(names []string) (id, name string, usage bytesize.ByteSize)
}, error) {
	var (
		ctx    = context.Background()
		result = calculateContainersMemoryUsageResult{}
	)

	containers, err := sa.dockerCli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get container list from docker")
	}

	result.containers = make([]struct {
		id, name string
		usage    bytesize.ByteSize
	}, len(containers))

	for i, container := range containers {
		var stats types.ContainerStats
		if stats, err = sa.dockerCli.ContainerStats(ctx, container.ID, false); err != nil {
			return nil, errors.Wrap(err, "failed to get container stats from docker")
		}

		v := &types.StatsJSON{}
		if err = json.NewDecoder(stats.Body).Decode(v); err != nil {
			return nil, errors.Wrap(err, "failed to decode stats response body to struct")
		}

		result.containers[i] = struct {
			id, name string
			usage    bytesize.ByteSize
		}{
			id: v.ID, name: v.Name,
			usage: getMemoryUsageSizeFrom(v),
		}
	}

	return nil, nil
}
