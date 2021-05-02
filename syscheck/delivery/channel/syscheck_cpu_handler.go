// Create file in v.1.0.0
// in syscheck_cpu_handler.go file, define delivery from channel msg to cpu usecase handler
// publishing msg to golang channel which is received from outside is not occurred in this package

package channel

import (
	"context"
	"log"
	"time"

	"github.com/DMS-SMS/v1-health-check/domain"
)

// cpuCheckHandler is delivered data handler about cpu check using usecase layer
type cpuCheckHandler struct {
	// handlerCtx is used for handling delivered channel using context
	handlerCtx systemCheckHandlerContext

	// cUsecase is usecase layer interface which is injected from package outside (maybe, in main)
	cUsecase domain.CPUCheckUseCase
}

// NewCPUCheckHandler define diskCheckHandler ptr instance & register handling channel msg to usecase
func NewCPUCheckHandler(c <-chan time.Time, cu domain.CPUCheckUseCase) {
	handler := &cpuCheckHandler{
		handlerCtx: globalContext,
		cUsecase:   cu,
	}

	go handler.startListening(c)
	log.Println("START TO LISTEN CHANNEL MSG ABOUT SYSTEM CPU CHECK")
}

// startListening method start listening using handlerCtx field startListening method
func (ch *cpuCheckHandler) startListening(c <-chan time.Time) {
	ch.handlerCtx.startListening(c, ch.checkCPU)
}

// checkCPU method set context & call usecase CheckCPU method, handle error
func (ch *cpuCheckHandler) checkCPU(t time.Time) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "time", t)

	if err := ch.cUsecase.CheckCPU(ctx); err != nil {
		log.Printf("error occurs in CheckCPU, err: %v", err)
	}
}
