// Create file in v.1.0.0
// agent_cpu_result.go is file that define result struct using as return value in agent method about cpu
// all result struct implement interface defined in return type of method signature in agent_cpu.go

package system

import "runtime"

// calculateContainersCPUUsageResult is result type of CalculateContainersCPUUsage
type calculateContainersCPUUsageResult struct {
	// cpuNum represent number of cpu core
	cpuNum int

	// containers is to keep cpu usage each of container get from GetTotalCPUUsage
	containers []struct {
		id, name   string
		cpuPercent float64
	}
}

// TotalCPUUsage return total cpu usage in docker containers
func (result calculateContainersCPUUsageResult) TotalCPUUsage() (usage float64) {
	var percent float64 = 0
	for _, container := range result.containers {
		percent += container.cpuPercent
	}
	usage = float64(runtime.NumCPU()) / 100 * percent
	return
}
