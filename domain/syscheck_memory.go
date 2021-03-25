// Create file in v.1.0.0
// syscheck_mem.go is file that declare model struct & repo interface about memory health check in syscheck domain.
// also, additional method of model struct is declared in this file, too.

package domain

import "github.com/inhies/go-bytesize"

// MemCheckHistory model is used for record memory health check history and result
type MemoryCheckHistory struct {
	// get required component by embedding systemCheckHistoryComponent
	systemCheckHistoryComponent

	// TotalUsageMemory specifies current memory usage of runtime system looked in memory check
	TotalUsageMemory bytesize.ByteSize

	// DockerUsageMemory specifies current total memory usage of docker looked in memory check when weak detected
	DockerUsageMemory float64

	// TemporaryFreeMemory specifies temporary freed memory size while recovering memory health
	TemporaryFreeMemory float64

	// MostMemoryConsumeContainer specifies the container name which is consumed most memory
	MostMemoryConsumeContainer string
}
