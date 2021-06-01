// Create file in v.1.0.0
// syscheck_mem.go is file that declare model struct & repo interface about memory health check in syscheck domain.
// also, additional method of model struct is declared in this file, too.

package domain

import (
	"context"
	"github.com/inhies/go-bytesize"
)

// MemCheckHistory model is used for record memory health check history and result
type MemoryCheckHistory struct {
	// get required component by embedding systemCheckHistoryComponent
	systemCheckHistoryComponent

	// TotalUsageMemory specifies current memory usage of runtime system looked in memory check
	TotalUsageMemory bytesize.ByteSize

	// DockerUsageMemory specifies current total memory usage of docker looked in memory check when weak detected
	DockerUsageMemory bytesize.ByteSize

	// TemporaryFreeMemory specifies temporary freed memory size while recovering memory health
	TemporaryFreeMemory bytesize.ByteSize

	// MostMemoryConsumeContainer specifies the container name which is consumed most memory
	MostMemoryConsumeContainer string
}

// MemoryCheckHistoryRepository is interface for repository layer used in usecase layer
// Repository is implemented with elasticsearch in v.1.0.0
type MemoryCheckHistoryRepository interface {
	// get required component by embedding systemCheckHistoryRepositoryComponent
	systemCheckHistoryRepositoryComponent

	// Store method save MemoryCheckHistory model in repository
	// b in return represents bytes of response body(map[string]interface{})
	Store(*MemoryCheckHistory) (b []byte, err error)
}

// MemoryCheckUseCase is interface used as business process handler about memory check
type MemoryCheckUseCase interface {
	// CheckMemory method check memory usage status and store memory check history using repository
	CheckMemory(ctx context.Context) error
}

// FillPrivateComponent overriding FillPrivateComponent method of systemCheckHistoryComponent
func (mc *MemoryCheckHistory) FillPrivateComponent() {
	mc.systemCheckHistoryComponent.FillPrivateComponent()
	mc._type = "MemoryCheck"
}

// DottedMapWithPrefix convert CPUCheckHistory to dotted map and return using MapWithPrefixKey of upper struct
// all key value of Map start with prefix received from parameter
func (mc *MemoryCheckHistory) DottedMapWithPrefix(prefix string) (m map[string]interface{}) {
	m = mc.systemCheckHistoryComponent.DottedMapWithPrefix(prefix)

	if prefix != "" {
		prefix += "."
	}

	// setting public field value in dotted map
	m[prefix+"total_usage_memory"] = mc.TotalUsageMemory.String()
	m[prefix+"docker_usage_memory"] = mc.DockerUsageMemory.String()
	m[prefix+"temporary_free_memory"] = mc.TemporaryFreeMemory.String()
	m[prefix+"most_memory_consume_container"] = mc.MostMemoryConsumeContainer

	return
}
