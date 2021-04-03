// Create file in v.1.0.0
// srvcheck_swarmpit_ucase.go is file that define usecase implementation about swarmpit check in srvcheck domain
// usecase layer depend on repository layer and is depended to delivery layer

package usecase

import (
	"github.com/docker/docker/api/types"
	"github.com/inhies/go-bytesize"
	"sync"

	"github.com/DMS-SMS/v1-health-check/domain"
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

	// dockerAgency is used as agency about docker engine API
	dockerAgency dockerAgency

	// status represent current process status of swarmpit health check
	status swarmpitCheckStatus

	// mutex help to prevent race condition when set status field value
	mutex sync.Mutex
}

// swarmpitCheckUsecaseConfig is the config getter interface for swarmpit check usecase
type swarmpitCheckUsecaseConfig interface {
	// get common config method from embedding serviceCheckUsecaseComponentConfig
	serviceCheckUsecaseComponentConfig

	// SwarmpitAppServiceName method returns string represent swarmpit app service name
	SwarmpitAppServiceName() string

	// JaegerIndexPattern method returns string represent jaeger index pattern
	SwarmpitAppMaxMemoryUsage() bytesize.ByteSize
}

// dockerAgency is agency that agent various command about docker engine API
type dockerAgency interface {
	// GetContainerWithServiceName return container which is instance of received service name
	GetContainerWithServiceName(srv string) (container interface {
		ID() string                     // get id of container
		Name() string                   // get name of container
		MemoryUsage() bytesize.ByteSize // get memory usage of container
	}, err error)

	// RemoveContainer remove container with id & option (auto created from docker swarm if exists)
	RemoveContainer(containerID string, options types.ContainerRemoveOptions) error
}

// NewSwarmpitCheckUsecase function return swarmpitCheckUsecase ptr instance after initializing
func NewSwarmpitCheckUsecase(
	cfg swarmpitCheckUsecaseConfig,
	shr domain.SwarmpitCheckHistoryRepository,
	sca slackChatAgency,
	da dockerAgency,
) domain.SwarmpitCheckUseCase {
	return &swarmpitCheckUsecase{
		// initialize field with parameter received from caller
		myCfg:           cfg,
		historyRepo:     shr,
		slackChatAgency: sca,
		dockerAgency:    da,

		// initialize field with default value
		status: swarmpitStatusHealthy,
		mutex:  sync.Mutex{},
	}
}
