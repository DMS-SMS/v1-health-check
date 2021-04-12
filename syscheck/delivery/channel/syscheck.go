// Create package in v.1.0.0
// delivery package is for delivery layer acted as presenter layer in syscheck domain which decide how the data will presented
// in delivery type, could be as REST API, gRPC, golang channel, or HTML file, etc ...
// in channel delivery, deliver data to usecase by receiving from golang channel while listening

// in this file, define global variable or function using in all struct defined in another file

package channel

import (
	"context"
)

var (
	// chanCancelCtx is used for checking if channel delivery is canceled by cancel method
	ChanCancelCtx context.Context

	// chanWaitGroup is used when start & end handling delivered chan with Add & Done method
	chanWaitGroup *sync.WaitGroup
)
