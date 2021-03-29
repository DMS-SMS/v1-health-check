// delivery package is for delivery layer acted as presenter layer in syscheck domain which decide how the data will presented
// in delivery type, could be as REST API, gRPC, golang channel, or HTML file, etc ...
// in channel delivery, deliver data to usecase by receiving from golang channel while listening

// in syscheck_disk_handler.go file, define delivery from channel msg to disk usecase handler
// publishing msg to golang channel which is received from outside is not occurred in this package

package channel

import (
	"context"
	"log"
	"time"

	"github.com/DMS-SMS/v1-health-check/domain"
)

// diskCheckHandler is delivered data handler about disk check using usecase layer
type diskCheckHandler struct {
	// DUsecase is usecase layer interface which is injected from package outside (maybe, in main)
	DUsecase domain.DiskCheckUseCase
}

// NewDiskCheckHandler define diskCheckHandler ptr instance & register handling channel msg to usecase
func NewDiskCheckHandler(c <-chan time.Time, du domain.DiskCheckUseCase) {
	handler := &diskCheckHandler{
		DUsecase: du,
	}

	handler.startListening(c)
}

// CheckDisk method set context & call usecase CheckDisk method, handle error
func (dh *diskCheckHandler) CheckDisk(t time.Time) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "time", t)

	if err := dh.DUsecase.CheckDisk(ctx); err != nil {
		log.Printf("error occurs in CheckDisk, err: %v", err)
	}
}

// startListening method start listening msg from golang channel & stream msg to another method
func (dh *diskCheckHandler) startListening(c <-chan time.Time) {
	for {
		select {
		case t := <-c:
			go dh.CheckDisk(t)
		}
	}
}
