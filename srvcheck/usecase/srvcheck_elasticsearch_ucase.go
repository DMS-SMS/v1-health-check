// Create file in v.1.0.0
// srvcheck_elasticsearch_ucase.go is file that define usecase implementation about elasticsearch check in srvcheck domain
// elasticsearch check usecase struct embed serviceCheckUsecaseComponent struct in ./srvcheck.go file

package usecase

import (
	"sync"
	"time"

	"github.com/DMS-SMS/v1-health-check/domain"
)

// elasticsearchCheckStatus is type to int constant represent current elasticsearch check process status
type elasticsearchCheckStatus int
const (
	elasticsearchStatusHealthy    elasticsearchCheckStatus = iota // represent elasticsearch check status is healthy
	elasticsearchStatusRecovering                                 // represent it's recovering elasticsearch status now
	elasticsearchStatusUnhealthy                                  // represent elasticsearch check status is unhealthy
)

// elasticsearchCheckUsecase implement ElasticsearchCheckUsecase interface in domain and used in delivery layer
type elasticsearchCheckUsecase struct {
	// myCfg is used for getting elasticsearch check usecase config
	myCfg elasticsearchCheckUsecaseConfig

	// historyRepo is used for store elasticsearch check history and injected from outside
	historyRepo domain.ElasticsearchCheckHistoryRepository

	// slackChat is used for agent slack API about chatting
	slackChatAgency slackChatAgency

	// elasticsearchAgency is used as agency about elasticsearch API
	elasticsearchAgency elasticsearchAgency

	// status represent current process status of elasticsearch health check
	status elasticsearchCheckStatus

	// mutex help to prevent race condition when set status field value
	mutex sync.Mutex
}

// elasticsearchCheckUsecaseConfig is the config getter interface for elasticsearch check usecase
type elasticsearchCheckUsecaseConfig interface {
	// get common config method from embedding serviceCheckUsecaseComponentConfig
	serviceCheckUsecaseComponentConfig

	// MaximumShardsNumber method returns int represent maximum shards number
	MaximumShardsNumber() int

	// JaegerIndexPattern method returns string represent jaeger index pattern
	JaegerIndexPattern() string

	// JaegerIndexLifeCycle method returns duration represent jaeger index life cycle
	JaegerIndexMinLifeCycle() time.Duration
}

// elasticsearchAgency is interface that agent elasticsearch with HTTP API
type elasticsearchAgency interface {
	// GetClusterHealth return interface have various get method about cluster health inform
	GetClusterHealth() (result interface{
		ActivePrimaryShards() int     // get active primary shards number in cluster health result
		ActiveShards() int            // get active shards number in cluster health result
		UnassignedShards() int        // get unassigned shards number in cluster health result
		ActiveShardsPercent() float64 // get active shards percent in cluster health result
	}, err error)

	// GetIndicesWithRegexp return indices list with regexp pattern
	GetIndicesWithPatterns(patterns []string) (indices []string, err error)

	// DeleteIndices method delete indices in list received from parameter
	DeleteIndices(indices []string) (err error)
}

// NewElasticsearchCheckUsecase function return elasticsearchCheckUseCase ptr instance after initializing
func NewElasticsearchCheckUsecase(
	cfg elasticsearchCheckUsecaseConfig,
	chr domain.ElasticsearchCheckHistoryRepository,
	sca slackChatAgency,
	ea elasticsearchAgency,
) domain.ElasticsearchCheckUseCase {
	return &elasticsearchCheckUsecase{
		// initialize field with parameter received from caller
		myCfg:               cfg,
		historyRepo:         chr,
		slackChatAgency:     sca,
		elasticsearchAgency: ea,

		// initialize field with default value
		status: elasticsearchStatusHealthy,
		mutex:  sync.Mutex{},
	}
}

// CheckElasticsearch check elasticsearch health with checkElasticsearch method & store check history in repository
// Implement CheckElasticsearch method of ElasticsearchCheckUseCase interface
func (ecu *elasticsearchCheckUsecase) CheckElasticsearch(ctx context.Context) (err error) {
	return
}

// method processed with below logic about elasticsearch health check according to current check status
// 0 : 정상적으로 인지된 상태 (상태 확인 수행)
// 0 -> 1 : Jaeger Index 삭제 실행 (Jaeger Index 삭제 알림 발행)
// 1 : Jaeger Index 삭제중 (상태 확인 수행 X)
// 1 -> 0 : Jaeger Index 삭제로 인해 상태 회복 완료 (상태 회복 알림 발행)
// 1 -> 2 : Jaeger Index 삭제를 해도 상태 회복 X (상태 회복 불가능 상태 알림 발행)
// 2 : 관리자가 직접 확인해야함 (상태 확인 수행 X)
// 2 -> 0 : 관리자 직접 상태 회복 완료 (상태 회복 알림 발행)
func (ecu *elasticsearchCheckUsecase) checkElasticsearch(ctx context.Context) (history *domain.ElasticsearchCheckHistory) {
	return
}
