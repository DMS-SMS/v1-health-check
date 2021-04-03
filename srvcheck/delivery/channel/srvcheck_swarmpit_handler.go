// in srvcheck_swarmpit_handler.go file, define delivery from channel msg to swarmpit check usecase handler
// publishing msg to golang channel which is received from outside is not occurred in this package

package channel

import (
	"context"
	"log"
	"time"

	"github.com/DMS-SMS/v1-health-check/domain"
)

// swarmpitCheckHandler is delivered data handler about swarmpit check using usecase layer
type swarmpitCheckHandler struct {
	// SUsecase is usecase layer interface which is injected from package outside (maybe, in main)
	SUsecase domain.SwarmpitCheckUseCase
}

// NewSwarmpitCheckHandler define swarmpitCheckHandler ptr instance & register handling channel msg to usecase
func NewSwarmpitCheckHandler(c <-chan time.Time, su domain.SwarmpitCheckUseCase) {
	handler := &swarmpitCheckHandler{
		SUsecase: su,
	}

	go handler.startListening(c)
	log.Println("START TO LISTEN CHANNEL MSG ABOUT SERVICE SWARMPIT CHECK")
}

// startListening method start listening msg from golang channel & stream msg to another method
func (sh *swarmpitCheckHandler) startListening(c <-chan time.Time) {
	for {
		select {
		case t := <-c:
			go sh.checkSwarmpit(t)
		}
	}
}

// checkSwarmpit method set context & call CheckSwarmpit usecase method, handle error
func (sh *swarmpitCheckHandler) checkSwarmpit(t time.Time) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "time", t)

	if err := sh.SUsecase.CheckSwarmpit(ctx); err != nil {
		log.Printf("error occurs in CheckSwarmpit, err: %v", err)
	}
}
