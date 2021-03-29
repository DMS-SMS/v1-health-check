// in syscheck_memory_handler.go file, define delivery from channel msg to memory usecase handler
// publishing msg to golang channel which is received from outside is not occurred in this package

package channel

import (

	"github.com/DMS-SMS/v1-health-check/domain"
)

// memoryCheckHandler is delivered data handler about memory check using usecase layer
type memoryCheckHandler struct {
	// CUsecase is usecase layer interface which is injected from package outside (maybe, in main)
	MUsecase domain.MemoryCheckUseCase
}

