// Create file in v.1.0.0
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
	// handlerCtx is used for handling delivered channel using context
	handlerCtx systemCheckHandlerContext

	// eUsecase is usecase layer interface which is injected from package outside (maybe, in main)
	eUsecase domain.ElasticsearchCheckUseCase
}

// NewElasticsearchCheckHandler define elasticsearchCheckHandler ptr instance & register handling channel msg to usecase
func NewElasticsearchCheckHandler(c <-chan time.Time, eu domain.ElasticsearchCheckUseCase) {
	handler := &elasticsearchCheckHandler{
		handlerCtx: globalContext,
		eUsecase:   eu,
	}

	go handler.startListening(c)
	log.Println("START TO LISTEN CHANNEL MSG ABOUT SERVICE ELASTICSEARCH CHECK")
}

// startListening method start listening using handlerCtx field startListening method
func (eh *elasticsearchCheckHandler) startListening(c <-chan time.Time) {
	eh.handlerCtx.startListening(c, eh.checkElasticsearch)
}

// checkElasticsearch method set context & call CheckElasticsearch usecase method, handle error
func (eh *elasticsearchCheckHandler) checkElasticsearch(t time.Time) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "time", t)

	if err := eh.eUsecase.CheckElasticsearch(ctx); err != nil {
		log.Printf("error occurs in CheckElasticsearch, err: %v", err)
	}
}
