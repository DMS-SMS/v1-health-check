// Create package in v.1.0.0
// Same as entities, struct and method in domain package will used in all layer.
// syscheck.go is file that define model as struct and abstract method of model as interface.
// Also, it declare usecase interface used as business layer.

// syscheck domain is managing the state of the system (CPU, memory, disk, etc.) periodically

// All model struct and interface is about 'system check' domain
// Most importantly, it only defines and does not implement interfaces.

package domain

import "time"

// SystemCheckHistory is basic model using by embedded in every model struct about check history
type SystemCheckHistory struct {
	// Agent specifies name of service that created this model
	Agent string

	// version specifies health checker version when this model was created
	Version string

	// Timestamp specifies the time at which this model was created.
	Timestamp time.Time
}
