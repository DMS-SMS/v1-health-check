// Create file in v.1.0.0
// syscheck_memory_ucase.go is file that define usecase implementation about syscheck memory domain
// memory check usecase struct embed systemCheckUsecaseComponent struct in ./syscheck.go file

package usecase

import (
	"github.com/inhies/go-bytesize"
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
	// myCfg is used for getting memory check usecase config
	myCfg memoryCheckUsecaseConfig

	// historyRepo is used for store memory check history and injected from outside
	historyRepo domain.MemoryCheckHistoryRepository

	// slackChat is used for agent slack API about chatting
	slackChatAgency slackChatAgency

	// memorySysAgency is used as agency about memory system command
	memorySysAgency memorySysAgency

	// dockerAgency is used as agency about docker command
	dockerAgency dockerAgency

	// status represent current process status of memory health check
	status memoryCheckStatus

	// mutex help to prevent race condition when set status field value
	mutex sync.Mutex
}

// memoryCheckUsecaseConfig is the config getter interface for memory check usecase
type memoryCheckUsecaseConfig interface {
	// get common config method from embedding systemCheckUsecaseComponentConfig
	systemCheckUsecaseComponentConfig

	// MemoryWarningUsage method returns float64 represent memory warning usage
	MemoryWarningUsage() float64

	// MemoryMaximumUsage method returns float64 represent memory maximum usage
	MemoryMaximumUsage() float64

	// MemoryMinimumUsageToRemove method returns float64 represent memory minimum usage to remove
	MemoryMinimumUsageToRemove() float64
}

// memorySysAgency is agency that agent various command about memory system
type memorySysAgency interface {
	// GetTotalSystemMemoryUsage return total memory usage as bytesize in system
	GetTotalSystemMemoryUsage() (size bytesize.ByteSize, err error)

	// CalculateContainersMemoryUsage calculate container memory usage & return result interface implementation
	CalculateContainersMemoryUsage() (result interface {
		// TotalMemoryUsage return total memory usage in docker containers
		TotalMemoryUsage() (size bytesize.ByteSize)

		// MostConsumerExceptFor return container consume the most memory except container names received from param
		MostConsumerExceptFor(names []string) (id, name string, size bytesize.ByteSize)
	}, err error)
}

// NewMemoryCheckUsecase function return memoryCheckUsecase ptr instance after initializing
func NewMemoryCheckUsecase(
	cfg memoryCheckUsecaseConfig,
	mhr domain.MemoryCheckHistoryRepository,
	sca slackChatAgency,
	msa memorySysAgency,
	da dockerAgency,
) domain.MemoryCheckUseCase {
	return &memoryCheckUsecase{
		// initialize field with parameter received from caller
		myCfg:           cfg,
		historyRepo:     mhr,
		slackChatAgency: sca,
		memorySysAgency: msa,
		dockerAgency:    da,

		// initialize field with default value
		status: memoryStatusHealthy,
		mutex:  sync.Mutex{},
	}
}
