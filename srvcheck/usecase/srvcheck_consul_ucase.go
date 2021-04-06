// Create file in v.1.0.0
// srvcheck_consul_ucase.go is file that define usecase implementation about consul check in srvcheck domain
// usecase layer depend on repository layer and is depended to delivery layer

package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
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

	// gRPCAgency is used as agency about gRPC
	gRPCAgency gRPCAgency

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
	ga gRPCAgency,
	da dockerAgency,
) domain.ConsulCheckUseCase {
	return &consulCheckUsecase{
		// initialize field with parameter received from caller
		myCfg:           cfg,
		historyRepo:     shr,
		slackChatAgency: sca,
		consulAgency:    ca,
		gRPCAgency:      ga,
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

// method processed with below logic about consul health check according to current check status
// 0 : 정상적으로 인지된 상태 (상태 확인 수행) (모든 등록된 Service 정상 작동 & 서비스별 인스턴스 최소 1개 존재)
// 0 -> 1 : Consul 상태 회복(작동X 노드 삭제 or 특정 서비스 재실행) 실행 (Consul 상태 회복 실행 알림 발행)
// 1 : Consul 상태 회복중 (상태 확인 수행 X)
// 1 -> 0 : Consul 상태 회복으로 인해 상태 회복 완료 (상태 회복 알림 발행)
// 1 -> 2 : Consul 상태 회복을 해도 상태 회복 X (상태 회복 불가능 상태 알림 발행)
// 2 : 관리자가 직접 확인해야함 (상태 확인 수행 X)
// 2 -> 0 : 관리자 직접 상태 회복 완료 (상태 회복 알림 발행)
func (ccu *consulCheckUsecase) checkConsul(ctx context.Context) (history *domain.ConsulCheckHistory) {
	_uuid := uuid.New().String()
	history = new(domain.ConsulCheckHistory)
	history.FillPrivateComponent()
	history.UUID = _uuid
	history.InstancesPerService = map[string][]string{}

	switch ccu.status {
	case consulStatusHealthy:
		break
	case consulStatusRecovering:
		history.ProcessLevel.Set(recoveringLevel)
		history.Message = "recovering consul health is already on process"
		return
	case consulStatusUnhealthy:
		history.ProcessLevel.Set(unhealthyLevel)
		history.Message = "consul check is unhealthy now"
		return
	}

	srvM := map[string][]struct{ id, addr string }{}
	for _, srv := range ccu.myCfg.CheckTargetServices() {
		cslSrv := ccu.myCfg.ConsulServiceNameSpace() + srv
		iter, err := ccu.consulAgency.GetServices(cslSrv)
		if err != nil {
			history.ProcessLevel.Set(errorLevel)
			history.SetError(errors.Wrap(err, "failed to get services"))
			return
		}

		if _, ok := history.InstancesPerService[cslSrv]; !ok {
			history.InstancesPerService[cslSrv] = []string{}
		}
		if _, ok := srvM[cslSrv]; !ok {
			srvM[cslSrv] = []struct{ id, addr string }{}
		}

		for iter.HasNext() {
			id, addr := iter.Next()
			history.InstancesPerService[cslSrv] = append(history.InstancesPerService[cslSrv], id)
			srvM[cslSrv] = append(srvM[cslSrv], struct{ id, addr string }{id: id, addr: addr})
		}
	}

	return
}

// setStatus set status field value using mutex Lock & Unlock
func (ccu *consulCheckUsecase) setStatus(status consulCheckStatus) {
	ccu.mutex.Lock()
	defer ccu.mutex.Unlock()
	ccu.status = status
}
