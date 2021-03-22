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

// GetTotalCPUUsage method return total using cpu usage with float64
func (sa *sysAgent) GetTotalCPUUsage() (usage float64, err error) {
	var (
		ctx = context.Background()
	)

	var lists []types.Container
	if lists, err = sa.dockerCli.ContainerList(ctx, types.ContainerListOptions{}); err != nil {
		err = errors.Wrap(err, "failed to get container list from docker")
		return
	}

	for _, list := range lists {
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
		usage += getCPUUsagePercentFrom(v)
	}

	usage = float64(runtime.NumCPU()) / 100 * usage
	return
}
