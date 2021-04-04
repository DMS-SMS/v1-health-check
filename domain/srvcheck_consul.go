// Create file in v.1.0.0
// srvcheck_consul.go is file that declare model struct & repo interface about consul check in srvcheck domain.
// also, additional method of model struct is declared in this file, too.

package domain

import "context"

// ConsulCheckHistory model is used for record consul check history and result
type ConsulCheckHistory struct {
	// get required component by embedding serviceCheckHistoryComponent
	serviceCheckHistoryComponent

	// InstancesPerService specifies instances list per service in consul
	InstancesPerService map[string][]string

	// IfInstanceDeregistered specifies if any instance in consul was deregistered
	IfInstanceDeregistered bool

	// DeregisteredInstances specifies id list of deregistered instance in consul check
	DeregisteredInstances []string

	// IfContainerRestarted specifies if container about service is restarted
	IfContainerRestarted bool

	// RestartedContainers specifies id list of restarted containers in consul check
	RestartedContainers []string
}

// ConsulCheckHistoryRepository is interface for repository layer used in usecase layer
// Repository is implemented with elasticsearch in v.1.0.0
type ConsulCheckHistoryRepository interface {
	// get required component by embedding serviceCheckHistoryRepositoryComponent
	serviceCheckHistoryRepositoryComponent

	// Store method save ConsulCheckHistory model in repository
	// b in return represents bytes of response body(map[string]interface{})
	Store(*ConsulCheckHistory) (b []byte, err error)
}

// ConsulCheckUseCase is interface used as business process handler about consul check
type ConsulCheckUseCase interface {
	// CheckConsul method check consul status and store check history using repository
	CheckConsul(ctx context.Context) error
}
