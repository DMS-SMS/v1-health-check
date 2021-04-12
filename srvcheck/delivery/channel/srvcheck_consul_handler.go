// Create file in v.1.0.0
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

// NewConsulCheckHandler define consulCheckHandler ptr instance & register handling channel msg to usecase
func NewConsulCheckHandler(c <-chan time.Time, cu domain.ConsulCheckUseCase) {
	handler := &consulCheckHandler{
		CUsecase: cu,
	}

	go handler.startListening(c)
	log.Println("START TO LISTEN CHANNEL MSG ABOUT SERVICE CONSUL CHECK")
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
