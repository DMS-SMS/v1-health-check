// Create file in v.1.0.0
// syscheck_memory_ucase.go is file that define usecase implementation about syscheck memory domain
// memory check usecase struct embed systemCheckUsecaseComponent struct in ./syscheck.go file

package usecase

import (
	"sync"

	"github.com/DMS-SMS/v1-health-check/domain"
)

// memoryCheckStatus is type to int constant represent current memory check process status
type memoryCheckStatus int
const (
	memoryStatusHealthy    memoryCheckStatus = iota // represent memory check status is healthy
	memoryStatusWarning                             // represent memory check status is warning now
	memoryStatusRecovering                          // represent it's recovering memory status now
	memoryStatusUnhealthy                           // represent memory check status is unhealthy
)

// memoryCheckUsecase implement MemoryCheckUsecase interface in domain and used in delivery layer
type memoryCheckUsecase struct {
	//// myCfg is used for getting memory check usecase config
	//myCfg memoryCheckUsecaseConfig

	// historyRepo is used for store memory check history and injected from outside
	historyRepo domain.MemoryCheckHistoryRepository

	// slackChat is used for agent slack API about chatting
	slackChatAgency slackChatAgency

	//// memorySysAgency is used as agency about memory system command
	//memorySysAgency memorySysAgency

	// dockerAgency is used as agency about docker command
	dockerAgency dockerAgency

	// status represent current process status of memory health check
	status memoryCheckStatus

	// mutex help to prevent race condition when set status field value
	mutex sync.Mutex
}
