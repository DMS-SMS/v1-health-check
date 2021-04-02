// delivery package is for delivery layer acted as presenter layer in srvcheck domain which decide how the data will presented
// in delivery type, could be as REST API, gRPC, golang channel, or HTML file, etc ...
// in channel delivery, deliver data to usecase by receiving from golang channel while listening

// in srvcheck_elasticsearch_handler.go file, define delivery from channel msg to elasticsearch check usecase handler
// publishing msg to golang channel which is received from outside is not occurred in this package

package channel

import (
	"github.com/DMS-SMS/v1-health-check/domain"
)

// elasticsearchCheckHandler is delivered data handler about elasticsearch check using usecase layer
type elasticsearchCheckHandler struct {
	// EUsecase is usecase layer interface which is injected from package outside (maybe, in main)
	EUsecase domain.ElasticsearchCheckUseCase
}
