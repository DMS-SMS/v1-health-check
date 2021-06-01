// Create file in v.1.0.0
// srvcheck_swarmpit.go is file that declare model struct & repo interface about swarmpit check in srvcheck domain.
// also, additional method of model struct is declared in this file, too.

package domain

import (
	"context"
	"github.com/inhies/go-bytesize"
)

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

// SwarmpitCheckUseCase is interface used as business process handler about swarmpit check
type SwarmpitCheckUseCase interface {
	// CheckSwarmpit method check swarmpit status and store check history using repository
	CheckSwarmpit(ctx context.Context) error
}

// FillPrivateComponent overriding FillPrivateComponent method of serviceCheckHistoryComponent
func (sh *SwarmpitCheckHistory) FillPrivateComponent() {
	sh.serviceCheckHistoryComponent.FillPrivateComponent()
	sh._type = "SwarmpitCheck"
}

// DottedMapWithPrefix convert SwarmpitCheckHistory to dotted map and return using MapWithPrefixKey of upper struct
// all key value of Map start with prefix received from parameter
func (sh *SwarmpitCheckHistory) DottedMapWithPrefix(prefix string) (m map[string]interface{}) {
	m = sh.serviceCheckHistoryComponent.DottedMapWithPrefix(prefix)

	if prefix != "" {
		prefix += "."
	}

	// setting public field value in dotted map
	m[prefix+"swarmpit_app_memory_usage"] = sh.SwarmpitAppMemoryUsage.String()
	m[prefix+"if_swarmpit_app_restarted"] = sh.IfSwarmpitAppRestarted

	return
}
