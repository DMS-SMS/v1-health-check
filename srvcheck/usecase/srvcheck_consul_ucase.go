// Create file in v.1.0.0
// srvcheck_consul_ucase.go is file that define usecase implementation about consul check in srvcheck domain
// usecase layer depend on repository layer and is depended to delivery layer

package usecase

import (
	"github.com/DMS-SMS/v1-health-check/domain"
	"sync"
	"time"
)

// consulCheckStatus is type to int constant represent current consul check process status
type consulCheckStatus int
const (
	consulStatusHealthy    consulCheckStatus = iota // represent consul check status is healthy
	consulStatusRecovering                          // represent it's recovering consul status now
	consulStatusUnhealthy                           // represent consul check status is unhealthy
)

// consulCheckUsecase implement ConsulCheckUsecase interface in domain and used in delivery layer
type consulCheckUsecase struct {
	// myCfg is used for getting consul check usecase config
	myCfg consulCheckUsecaseConfig

	// historyRepo is used for store consul check history and injected from outside
	historyRepo domain.ConsulCheckHistoryRepository

	// slackChat is used for agent slack API about chatting
	slackChatAgency slackChatAgency

	// status represent current process status of consul health check
	status consulCheckStatus

	// mutex help to prevent race condition when set status field value
	mutex sync.Mutex
}

// consulCheckUsecaseConfig is the config getter interface for consul check usecase
type consulCheckUsecaseConfig interface {
	// get common config method from embedding serviceCheckUsecaseComponentConfig
	serviceCheckUsecaseComponentConfig

	// CheckTargetServices method returns string slice containing target services to check in usecase
	CheckTargetServices() []string

	// HealthCheckPingTimeOut method returns bytesize represent swarmpit app maximum memory usage
	HealthCheckPingTimeOut() time.Duration
}
