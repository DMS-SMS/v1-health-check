// Create file in v.1.0.0
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
	// handlerCtx is used for handling delivered channel using context
	handlerCtx systemCheckHandlerContext

	// mUsecase is usecase layer interface which is injected from package outside (maybe, in main)
	mUsecase domain.MemoryCheckUseCase
}

// NewDiskCheckHandler define memoryCheckHandler ptr instance & register handling channel msg to usecase
func NewMemoryCheckHandler(c <-chan time.Time, mu domain.MemoryCheckUseCase) {
	handler := &memoryCheckHandler{
		handlerCtx: globalContext,
		mUsecase:   mu,
	}

	go handler.startListening(c)
	log.Println("START TO LISTEN CHANNEL MSG ABOUT SYSTEM MEMORY CHECK")
}

// startListening method start listening using handlerCtx field startListening method
func (mh *memoryCheckHandler) startListening(c <-chan time.Time) {
	mh.handlerCtx.startListening(c, mh.checkMemory)
}

// checkMemory method set context & call usecase CheckMemory method, handle error
func (mh *memoryCheckHandler) checkMemory(t time.Time) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "time", t)

	if err := mh.mUsecase.CheckMemory(ctx); err != nil {
		log.Printf("error occurs in CheckMemory, err: %v", err)
	}
}
