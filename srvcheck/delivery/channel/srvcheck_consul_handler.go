// delivery package is for delivery layer acted as presenter layer in srvcheck domain which decide how the data will presented
// in delivery type, could be as REST API, gRPC, golang channel, or HTML file, etc ...
// in channel delivery, deliver data to usecase by receiving from golang channel while listening

// in srvcheck_consul_handler.go file, define delivery from channel msg to consul check usecase handler
// publishing msg to golang channel which is received from outside is not occurred in this package

package channel

import (
	"context"
	"log"
	"time"

	"github.com/DMS-SMS/v1-health-check/domain"
)

// consulCheckHandler is delivered data handler about consul check using usecase layer
type consulCheckHandler struct {
	// CUsecase is usecase layer interface which is injected from package outside (maybe, in main)
	CUsecase domain.ConsulCheckUseCase
}

// startListening method start listening msg from golang channel & stream msg to another method
func (ch *consulCheckHandler) startListening(c <-chan time.Time) {
	for {
		select {
		case t := <-c:
			go ch.checkConsul(t)
		}
	}
}

// checkConsul method set context & call CheckConsul usecase method, handle error
func (ch *consulCheckHandler) checkConsul(t time.Time) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "time", t)

	if err := ch.CUsecase.CheckConsul(ctx); err != nil {
		log.Printf("error occurs in CheckConsul, err: %v", err)
	}
}
