// delivery package is for delivery layer acted as presenter layer in srvcheck domain which decide how the data will presented
// in delivery type, could be as REST API, gRPC, golang channel, or HTML file, etc ...
// in channel delivery, deliver data to usecase by receiving from golang channel while listening

// in srvcheck_elasticsearch_handler.go file, define delivery from channel msg to elasticsearch check usecase handler
// publishing msg to golang channel which is received from outside is not occurred in this package

package channel

import (
	"context"
	"log"
	"time"

	"github.com/DMS-SMS/v1-health-check/domain"
)

// elasticsearchCheckHandler is delivered data handler about elasticsearch check using usecase layer
type elasticsearchCheckHandler struct {
	// EUsecase is usecase layer interface which is injected from package outside (maybe, in main)
	EUsecase domain.ElasticsearchCheckUseCase
}

// NewElasticsearchCheckHandler define elasticsearchCheckHandler ptr instance & register handling channel msg to usecase
func NewElasticsearchCheckHandler(c <-chan time.Time, eu domain.ElasticsearchCheckUseCase) {
	handler := &elasticsearchCheckHandler{
		EUsecase: eu,
	}

	go handler.startListening(c)
	log.Println("START TO LISTEN CHANNEL MSG ABOUT SERVICE ELASTICSEARCH CHECK")
}

// startListening method start listening msg from golang channel & stream msg to another method
func (eh *elasticsearchCheckHandler) startListening(c <-chan time.Time) {
	for {
		select {
		case t := <-c:
			go eh.checkElasticsearch(t)
		}
	}
}

// checkElasticsearch method set context & call CheckElasticsearch usecase method, handle error
func (eh *elasticsearchCheckHandler) checkElasticsearch(t time.Time) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "time", t)

	if err := eh.EUsecase.CheckElasticsearch(ctx); err != nil {
		log.Printf("error occurs in CheckElasticsearch, err: %v", err)
	}
}
