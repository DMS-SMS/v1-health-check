// Create file in v.1.0.0
// syscheck_disk_ucase.go is file that define usecase implementation about disk check domain
// disk check usecase struct embed systemCheckUsecaseComponent struct in ./syscheck.go file

package usecase

import "github.com/DMS-SMS/v1-health-check/domain"

// diskCheckUsecase implement DiskCheckUsecase interface in domain and used in delivery layer
type diskCheckUsecase struct {
	// get required component by embedding systemCheckUsecaseComponent struct
	systemCheckUsecaseComponent

	// historyRepo is used for store disk check history and injected from outside
	historyRepo domain.DiskCheckHistoryRepository
}

// diskCheckUsecaseConfig is the config getter interface for disk check usecase
type diskCheckUsecaseConfig interface {
	// get common config method from embedding systemCheckUsecaseComponentConfig
	systemCheckUsecaseComponentConfig

	// DiskMinCapacity method returns byte size represent disk minimum capacity
	DiskMinCapacity() bytesize.ByteSize
}
