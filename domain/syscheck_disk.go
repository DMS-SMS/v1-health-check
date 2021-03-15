// Create file in v.1.0.0
// syscheck_disk.go is file that declare model struct & repo interface about disk in syscheck domain.
// also, additional method of model struct is declared in this file, too.

package domain

import (
	"context"
	"github.com/inhies/go-bytesize"
)

// DiskCheckHistory model is used for record disk health check history and result
type DiskCheckHistory struct {
	// get required component by embedding systemCheckHistoryComponent
	systemCheckHistoryComponent

	// DiskCapacity specifies remain disk capacity of runtime system looked in disk check
	DiskCapacity bytesize.ByteSize
}

// DiskCheckHistoryRepository is abstract method used in business layer
// Repository is implemented with elastic search in v.1.0.0
type DiskCheckHistoryRepository interface {
	// get required component by embedding systemCheckHistoryRepositoryComponent
	systemCheckHistoryRepositoryComponent

	// Store method save DiskCheckHistory model in repository
	// b in return represents bytes of response body(map[string]interface{})
	Store(*DiskCheckHistory) (b []byte, err error)
}

// DiskCheckUseCase is interface used as business process handler about disk check
type DiskCheckUseCase interface {
	// CheckDisk method check disk capacity status and store disk check history using repository
	CheckDisk(ctx context.Context) error
}

// FillPrivateComponent overriding FillPrivateComponent method of systemCheckHistoryComponent
func (dh *DiskCheckHistory) FillPrivateComponent() {
	dh.systemCheckHistoryComponent.FillPrivateComponent()
}

// DottedMapWithPrefix convert DiskCheckHistory to dotted map and return using MapWithPrefixKey of upper struct
// all key value of Map start with prefix received from parameter
func (dh *DiskCheckHistory) DottedMapWithPrefix(prefix string) (m map[string]interface{}) {
	m = dh.systemCheckHistoryComponent.DottedMapWithPrefix(prefix)

	if prefix != "" {
		prefix += "."
	}

	// setting public field value in dotted map
	m[prefix + "disk_capacity"] = dh.DiskCapacity.String()

	return
}
