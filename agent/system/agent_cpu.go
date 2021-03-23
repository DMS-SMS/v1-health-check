// Create file in v.1.0.0
// agent_cpu.go is file that define method of sysAgent that agent command about cpu
// For example in cpu command, there are get total cpu usage, prune cpu, etc ...

package system

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
	"runtime"
)

// CalculateContainersCPUUsage calculate cpu usage & return calculateContainersCPUUsageResult
func (sa *sysAgent) CalculateContainersCPUUsage() (interface {
	TotalCPUUsage() (usage float64)
	MostConsumerExceptFor([]string) (id, name string, usage float64)
}, error) {
	var (
		ctx    = context.Background()
		result = calculateContainersCPUUsageResult{}
	)

	lists, err := sa.dockerCli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get container list from docker")
	}

	result.cpuNum = runtime.NumCPU()
	result.containers = make([]struct {
		id, name string
		usage    float64
	}, len(lists))

	for i, list := range lists {
		var stats types.ContainerStats
		if stats, err = sa.dockerCli.ContainerStats(ctx, list.ID, false); err != nil {
			return nil, errors.Wrap(err, "failed to get container stats from docker")
		}

		v := &types.StatsJSON{}
		if err = json.NewDecoder(stats.Body).Decode(v); err != nil {
			return nil, errors.Wrap(err, "failed to decode stats response body to struct")
		}

		result.containers[i] = struct {
			id, name string
			usage    float64
		}{
			id: v.ID, name: v.Name,
			usage: float64(result.cpuNum) / 100 * getCPUUsagePercentFrom(v),
		}
	}

	return result, nil
}

// getCPUUsagePercentFrom get cpu usage as percent from types.StatsJson struct
func getCPUUsagePercentFrom(v *types.StatsJSON) (per float64) {
	// calculate the change for the cpu usage of the container in between readings
	cpuDelta := float64(v.CPUStats.CPUUsage.TotalUsage) - float64(v.PreCPUStats.CPUUsage.TotalUsage)
	// calculate the change for the entire system between readings
	systemDelta := float64(v.CPUStats.SystemUsage) - float64(v.PreCPUStats.SystemUsage)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		per = (cpuDelta / systemDelta) * float64(len(v.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}

	return
}

//func is
