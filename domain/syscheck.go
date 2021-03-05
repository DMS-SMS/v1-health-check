// Create package in v.1.0.0
// Same as entities, struct and method in domain package will used in all layer.
// syscheck.go is file that define model as struct and abstract method of model as interface.
// Also, it declare usecase interface used as business layer.

// syscheck domain is managing the state of the system (CPU, memory, disk, etc.) periodically

// All model struct and interface is about 'system check' domain
// Most importantly, it only defines and does not implement interfaces.

package domain

import "time"

// systemCheckHistoryComponent is basic model using by embedded in every model struct about check history
type systemCheckHistoryComponent struct {
	// Agent specifies name of service that created this model
	Agent string

	// version specifies health checker version when this model was created
	Version string

	// Timestamp specifies the time when this model was created.
	Timestamp time.Time
}

// systemCheckHistoryRepositoryComponent is basic interface using by embedded in every repository about check history
type systemCheckHistoryRepositoryComponent interface {
	// Migrate method build environment for storage in stores such as Mysql or Elasticsearch, etc.
	Migrate() error
}

// DiskCheckHistory model is used for record disk health check history and result
type DiskCheckHistory struct {
	// get required component by embedding systemCheckHistoryComponent
	systemCheckHistoryComponent
}

// DiskCheckHistoryRepository is abstract method used in business layer
// Repository is implemented with elastic search in v.1.0.0
type DiskCheckHistoryRepository interface {
	// get required component by embedding systemCheckHistoryRepositoryComponent
	systemCheckHistoryRepositoryComponent

	// Store method save DiskCheckHistory model in repository
	Store(*DiskCheckHistory) error
}
