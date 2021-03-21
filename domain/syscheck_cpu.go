// Create file in v.1.0.0
// syscheck_cpu.go is file that declare model struct & repo interface about cpu health check in syscheck domain.
// also, additional method of model struct is declared in this file, too.

package domain

import (
	"context"
	"github.com/inhies/go-bytesize"
)

// CPUCheckHistory model is used for record cpu health check history and result
type CPUCheckHistory struct {
	// get required component by embedding systemCheckHistoryComponent
	systemCheckHistoryComponent

	// Usage specifies current cpu usage of runtime system looked in cpu check
	Usage bytesize.ByteSize

	// Free specifies freed cpu size while recovering cpu health
	Free bytesize.ByteSize
}

// CPUCheckHistoryRepository is interface for repository layer used in usecase layer
// Repository is implemented with elasticsearch in v.1.0.0
type CPUCheckHistoryRepository interface {
	// get required component by embedding systemCheckHistoryRepositoryComponent
	systemCheckHistoryRepositoryComponent

	// Store method save CPUCheckHistory model in repository
	// b in return represents bytes of response body(map[string]interface{})
	Store(*CPUCheckHistory) (b []byte, err error)
}

// DiskCheckUseCase is interface used as business process handler about cpu check
type CPUCheckUseCase interface {
	// CheckCPU method check cpu usage status and store cpu check history using repository
	CheckCPU(ctx context.Context) error
}

// FillPrivateComponent overriding FillPrivateComponent method of systemCheckHistoryComponent
func (ch *CPUCheckHistory) FillPrivateComponent() {
	ch.systemCheckHistoryComponent.FillPrivateComponent()
	ch._type = "CPUCheck"
}
