// Create file in v.1.0.0
// agent_cpu_result.go is file that define result struct using as return value in agent method about cpu
// all result struct implement interface defined in return type of method signature in agent_cpu.go

package system

import (
	"sort"
	"strings"
)

// calculateContainersCPUUsageResult is result type of CalculateContainersCPUUsage
type calculateContainersCPUUsageResult struct {
	// cpuNum represent number of cpu core
	cpuNum int

	// containers is to keep cpu usage each of container get from GetTotalCPUUsage
	containers []struct {
		id, name   string
		usage float64
	}
}

// TotalCPUUsage return total cpu usage in docker containers
func (result calculateContainersCPUUsageResult) TotalCPUUsage() (usage float64) {
	for _, container := range result.containers {
		usage += container.usage
	}
	return
}

// MostConsumerExceptFor handle logic using mostConsumerExceptFor method
func (result calculateContainersCPUUsageResult) MostConsumerExceptFor(excepts []string) (id, name string, usage float64) {
	m := map[string]bool{}
	for _, except := range excepts {
		m[except] = true
	}
	return result.mostConsumerExceptFor(m)
}

// mostConsumerExceptFor return most CPU consumer inform except for container names received from param
func (result calculateContainersCPUUsageResult) mostConsumerExceptFor(excepts map[string]bool) (id, name string, usage float64) {
	sort.Slice(result.containers, func(i, j int) bool {
		return result.containers[i].usage > result.containers[j].usage
	})

	for _, container := range result.containers {
		sep := strings.Split(strings.TrimPrefix(container.name, "/"), ".")
		if len(sep) == 0 {
			continue
		}
		if _, ok := excepts[sep[0]]; ok {
			continue
		}

		id = container.id
		name = container.name
		usage = container.usage
		break
	}

	return
}
