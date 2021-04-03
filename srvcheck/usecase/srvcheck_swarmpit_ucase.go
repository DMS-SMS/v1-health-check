// Create file in v.1.0.0
// srvcheck_swarmpit_ucase.go is file that define usecase implementation about swarmpit check in srvcheck domain
// usecase layer depend on repository layer and is depended to delivery layer

package usecase

import (
	"github.com/DMS-SMS/v1-health-check/domain"
	"sync"
)

// swarmpitCheckStatus is type to int constant represent current swarmpit check process status
type swarmpitCheckStatus int
const (
	swarmpitStatusHealthy    swarmpitCheckStatus = iota // represent swarmpit check status is healthy
	swarmpitStatusRecovering                            // represent it's recovering swarmpit status now
	swarmpitStatusUnhealthy                             // represent swarmpit check status is unhealthy
)

// swarmpitCheckUsecase implement SwarmpitCheckUsecase interface in domain and used in delivery layer
type swarmpitCheckUsecase struct {
	// myCfg is used for getting swarmpit check usecase config
	myCfg swarmpitCheckUsecaseConfig

	// historyRepo is used for store swarmpit check history and injected from outside
	historyRepo domain.SwarmpitCheckHistoryRepository

	// slackChat is used for agent slack API about chatting
	slackChatAgency slackChatAgency

	// status represent current process status of swarmpit health check
	status swarmpitCheckStatus

	// mutex help to prevent race condition when set status field value
	mutex sync.Mutex
}
