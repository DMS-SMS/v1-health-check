// Create file in v.1.0.0
// srvcheck_swarmpit.go is file that declare model struct & repo interface about swarmpit check in srvcheck domain.
// also, additional method of model struct is declared in this file, too.

package domain

import "github.com/inhies/go-bytesize"

// SwarmpitCheckHistory model is used for record swarmpit check history and result
type SwarmpitCheckHistory struct {
	// get required component by embedding serviceCheckHistoryComponent
	serviceCheckHistoryComponent

	// SwarmpitAppMemoryUsage specifies memory usage of swarmpit app container
	SwarmpitAppMemoryUsage bytesize.ByteSize

	// IfSwarmpitAppRestarted specifies if swarmpit container was restarted
	IfSwarmpitAppRestarted bool
}

// SwarmpitCheckHistoryRepository is interface for repository layer used in usecase layer
// Repository is implemented with elasticsearch in v.1.0.0
type SwarmpitCheckHistoryRepository interface {
	// get required component by embedding serviceCheckHistoryRepositoryComponent
	serviceCheckHistoryRepositoryComponent

	// Store method save SwarmpitCheckHistory model in repository
	// b in return represents bytes of response body(map[string]interface{})
	Store(*SwarmpitCheckHistory) (b []byte, err error)
}
