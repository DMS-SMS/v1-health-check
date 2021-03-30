// Create package in v.1.0.0
// Same as entities, struct and method in domain package will used in all layer.
// srvcheck.go is file that define model as struct and abstract method of model as interface.
// Also, it declare usecase interface used as business layer.

// srvcheck domain is managing the state of the service (elasticsearch, swarm, consul, etc ...) periodically

// All model struct and interface is about service check domain
// Most importantly, it only defines and does not implement interfaces.

package domain

import "time"

// serviceCheckHistoryComponent is basic model using by embedded in every model struct about service check history
type serviceCheckHistoryComponent struct {
	// private field in below, these fields have fixed value so always set in FillPrivateComponent method
	// Agent specifies name of service that created this model
	agent string

	// version specifies health checker version when this model was created
	version string

	// Timestamp specifies the time when this model was created.
	timestamp time.Time

	// Domain specifies domain about right this package, srvcheck
	domain string

	// _type specifies detail service type in service check domain (Ex, elasticsearch, swarm, consul, etc ...)
	_type string
}
