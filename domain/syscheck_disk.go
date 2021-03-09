// Create file in v.1.0.0
// syscheck_disk.go is file that declare model struct & repo interface about disk in syscheck domain.
// also, additional method of model struct is declared in this file, too.

package domain

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

