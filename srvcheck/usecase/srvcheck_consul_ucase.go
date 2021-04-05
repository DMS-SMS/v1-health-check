// Create file in v.1.0.0
// srvcheck_consul_ucase.go is file that define usecase implementation about consul check in srvcheck domain
// usecase layer depend on repository layer and is depended to delivery layer

package usecase

import (
	"context"
	"github.com/pkg/errors"
	"sync"
	"time"

	"github.com/DMS-SMS/v1-health-check/domain"
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

	// consulAgency is used as agency about consul API
	consulAgency consulAgency

	// dockerAgency is used as agency about docker engine API
	dockerAgency dockerAgency

	// status represent current process status of consul health check
	status consulCheckStatus

	// mutex help to prevent race condition when set status field value
	mutex sync.Mutex
}

// consulAgency is agency that agent various command about consul API
type consulAgency interface {
	// GetServices method get services in consul & return services interface implement
	GetServices(srv string) (srvIter interface {
		HasNext() bool           // HasNext method return if srvIter has next element
		Next() (id, addr string) // Next method return next service id, address
	}, err error)

	// DeregisterInstance method deregister service in consul with received id
	DeregisterInstance(id string) (err error)
}

// gRPCAgency is agency that agent various command about gRPC
type gRPCAgency interface {
	// PingToCheckConn ping for connection check to gRPC node
	PingToCheckConn(ctx context.Context, target string, opts ...grpc.DialOption) (err error)
}

// consulCheckUsecaseConfig is the config getter interface for consul check usecase
type consulCheckUsecaseConfig interface {
	// get common config method from embedding serviceCheckUsecaseComponentConfig
	serviceCheckUsecaseComponentConfig

	// CheckTargetServices method returns string slice containing target services to check in usecase
	CheckTargetServices() []string

	// ConsulServiceNameSpace method returns name space of consul service
	ConsulServiceNameSpace() string

	// DockerServiceNameSpace method returns name space of docker service
	DockerServiceNameSpace() string

	// ConnCheckPingTimeOut method returns timeout duration in ping to check connection
	ConnCheckPingTimeOut() time.Duration
}

// NewConsulCheckUsecase function return ConsulCheckUseCase implementation after initializing
func NewConsulCheckUsecase(
	cfg consulCheckUsecaseConfig,
	shr domain.ConsulCheckHistoryRepository,
	sca slackChatAgency,
	ca consulAgency,
	da dockerAgency,
) domain.ConsulCheckUseCase {
	return &consulCheckUsecase{
		// initialize field with parameter received from caller
		myCfg:           cfg,
		historyRepo:     shr,
		slackChatAgency: sca,
		consulAgency:    ca,
		dockerAgency:    da,

		// initialize field with default value
		status: consulStatusHealthy,
		mutex:  sync.Mutex{},
	}
}

// CheckConsul check consul health with checkConsul method & store check history in repository
// Implement CheckConsul method of ConsulCheckUseCase interface
func (ccu *consulCheckUsecase) CheckConsul(ctx context.Context) (err error) {
	history := ccu.checkConsul(ctx)

	if b, err := ccu.historyRepo.Store(history); err != nil {
		return errors.Wrapf(err, "failed to store consul check history, response: %s", string(b))
	}

	return
}
