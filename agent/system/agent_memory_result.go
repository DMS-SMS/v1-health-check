// Create file in v.1.0.0
// agent_memory_result.go is file that define result struct using as return value in agent method about memory
// all result struct implement interface defined in return type of method signature in agent_memory.go

package system

import (
	"github.com/inhies/go-bytesize"
	"sort"
	"strings"
)

// calculateContainersMemoryUsageResult is result type of CalculateContainersMemoryUsage
type calculateContainersMemoryUsageResult struct {
	// containers is to keep memory usage each of container get from CalculateContainersMemoryUsage
	containers []struct {
		id, name string
		usage    bytesize.ByteSize
	}
}

// TotalCPUUsage return total memory usage in docker containers
func (result calculateContainersMemoryUsageResult) TotalMemoryUsage() (usage bytesize.ByteSize) {
	for _, container := range result.containers {
		usage += container.usage
	}
	return
}

// MostConsumerExceptFor handle logic using mostConsumerExceptFor method
func (result calculateContainersMemoryUsageResult) MostConsumerExceptFor(excepts []string) (id, name string, usage bytesize.ByteSize) {
	m := map[string]bool{}
	for _, except := range excepts {
		m[except] = true
	}
	return result.mostConsumerExceptFor(m)
}

// mostConsumerExceptFor return most memory consumer inform except for container names received from param
func (result calculateContainersMemoryUsageResult) mostConsumerExceptFor(excepts map[string]bool) (id, name string, usage bytesize.ByteSize) {
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
