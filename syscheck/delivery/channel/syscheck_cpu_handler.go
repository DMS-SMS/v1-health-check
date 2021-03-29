// delivery package is for delivery layer acted as presenter layer in syscheck domain which decide how the data will presented
// in delivery type, could be as REST API, gRPC, golang channel, or HTML file, etc ...
// in channel delivery, deliver data to usecase by receiving from golang channel while listening

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
	// CUsecase is usecase layer interface which is injected from package outside (maybe, in main)
	CUsecase domain.CPUCheckUseCase
}

// NewDiskCheckHandler define diskCheckHandler ptr instance & register handling channel msg to usecase
func NewCPUCheckHandler(c <-chan time.Time, cu domain.CPUCheckUseCase) {
	handler := &cpuCheckHandler{
		CUsecase: cu,
	}

	handler.startListening(c)
}

// CheckCPU method set context & call usecase CheckCPU method, handle error
func (ch *cpuCheckHandler) CheckCPU(t time.Time) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "time", t)

	if err := ch.CUsecase.CheckCPU(ctx); err != nil {
		log.Printf("error occurs in CheckCPU, err: %v", err)
	}
}

// startListening method start listening msg from golang channel & stream msg to another method
func (ch *cpuCheckHandler) startListening(c <-chan time.Time) {
	for {
		select {
		case t := <-c:
			go ch.CheckCPU(t)
		}
	}
}
