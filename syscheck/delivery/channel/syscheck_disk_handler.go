// delivery package is for delivery layer acted as presenter layer in syscheck domain which decide how the data will presented
// in delivery type, could be as REST API, gRPC, golang channel, or HTML file, etc ...
// in channel delivery, deliver data to usecase by receiving from golang channel while listening

// in syscheck_disk_handler.go file, define delivery from channel msg to disk usecase handler
// publishing msg to golang channel which is received from outside is not occurred in this package

package channel

import (
	"github.com/DMS-SMS/v1-health-check/domain"
)

// diskCheckHandler is delivered data handler about disk check using usecase layer
type diskCheckHandler struct {
	// DUsecase is usecase layer interface which is injected from package outside (maybe, in main)
	DUsecase domain.DiskCheckUseCase
}
