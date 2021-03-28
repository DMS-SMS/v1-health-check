// Create file in v.1.0.0
// agent_memory_result.go is file that define result struct using as return value in agent method about memory
// all result struct implement interface defined in return type of method signature in agent_memory.go

package system

import "github.com/inhies/go-bytesize"

// calculateContainersMemoryUsageResult is result type of CalculateContainersMemoryUsage
type calculateContainersMemoryUsageResult struct {
	// containers is to keep memory usage each of container get from CalculateContainersMemoryUsage
	containers []struct {
		id, name string
		usage    bytesize.ByteSize
	}
}

// TotalCPUUsage return total memory usage in docker containers
func (result calculateContainersMemoryUsageResult) TotalCPUUsage() (usage bytesize.ByteSize) {
	for _, container := range result.containers {
		usage += container.usage
	}
	return
}
