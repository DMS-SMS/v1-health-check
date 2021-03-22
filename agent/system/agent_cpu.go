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

// calculateContainersCPUUsageResult is result type of CalculateContainersCPUUsage
type calculateContainersCPUUsageResult struct {
	// cpuNum represent number of cpu core
	cpuNum int

	// containers is to keep cpu usage each of container get from GetTotalCPUUsage
	containers []struct {
		id, name string
		cpuUsage float64
	}
}

// TotalCPUUsage return total cpu usage in docker containers
func (cr calculateContainersCPUUsageResult) TotalCPUUsage() (usage float64) {
	for _, container := range cr.containers {
		usage += container.cpuUsage
	}
	return
}

// CalculateContainersCPUUsage calculate cpu usage & return calculateContainersCPUUsageResult
func (sa *sysAgent) CalculateContainersCPUUsage() (usage float64, err error) {
	var (
		ctx = context.Background()
	)

	var lists []types.Container
	if lists, err = sa.dockerCli.ContainerList(ctx, types.ContainerListOptions{}); err != nil {
		err = errors.Wrap(err, "failed to get container list from docker")
		return
	}
	sa.containersCPUUsage = make([]struct{
		id, name string
		cpuUsage float64
	}, len(lists))

	for i, list := range lists {
		var stats types.ContainerStats
		if stats, err = sa.dockerCli.ContainerStats(ctx, list.ID, false); err != nil {
			err = errors.Wrap(err, "failed to get container stats from docker")
			return
		}

		v := &types.StatsJSON{}
		if err = json.NewDecoder(stats.Body).Decode(v); err != nil {
			err = errors.Wrap(err, "failed to decode stats response body to struct")
			return
		}

		sa.containersCPUUsage[i] = struct {
			id, name string
			cpuUsage float64
		}{
			id: v.ID, name: v.Name,
			cpuUsage: getCPUUsagePercentFrom(v),
		}
		usage += sa.containersCPUUsage[i].cpuUsage
	}

	usage = float64(runtime.NumCPU()) / 100 * usage
	return
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
