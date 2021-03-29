// in syscheck_memory_handler.go file, define delivery from channel msg to memory usecase handler
// publishing msg to golang channel which is received from outside is not occurred in this package

package channel

import (
	"context"
	"log"
	"time"

	"github.com/DMS-SMS/v1-health-check/domain"
)

// memoryCheckHandler is delivered data handler about memory check using usecase layer
type memoryCheckHandler struct {
	// CUsecase is usecase layer interface which is injected from package outside (maybe, in main)
	MUsecase domain.MemoryCheckUseCase
}

// NewDiskCheckHandler define memoryCheckHandler ptr instance & register handling channel msg to usecase
func NewMemoryCheckHandler(c <-chan time.Time, mu domain.MemoryCheckUseCase) {
	handler := &memoryCheckHandler{
		MUsecase: mu,
	}

	handler.startListening(c)
}

// CheckMemory method set context & call usecase CheckMemory method, handle error
func (mh *memoryCheckHandler) CheckMemory(t time.Time) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "time", t)

	if err := mh.MUsecase.CheckMemory(ctx); err != nil {
		log.Printf("error occurs in CheckMemory, err: %v", err)
	}
}

// startListening method start listening msg from golang channel & stream msg to another method
func (mh *memoryCheckHandler) startListening(c <-chan time.Time) {
	for {
		select {
		case t := <-c:
			mh.CheckMemory(t)
		}
	}
}
