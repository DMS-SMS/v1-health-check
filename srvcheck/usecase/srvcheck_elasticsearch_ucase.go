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
	elasticsearchStatusWarning                                    // represent elasticsearch check status is warning now
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

	// TargetIndices method returns string represent target indices separated with dot
	TargetIndices() []string

	// MaximumShardsNumber method returns int represent maximum shards number
	MaximumShardsNumber() int

	// JaegerIndexRegexp method returns string represent jaeger index regexp
	JaegerIndexRegexp() string

	// JaegerIndexLifeCycle method returns duration represent jaeger index life cycle
	JaegerIndexMinLifeCycle() time.Duration
}

// elasticsearchAgency is interface that agent elasticsearch with HTTP API
type elasticsearchAgency interface {
	// GetClusterHealth return interface have various get method about cluster health inform
	GetClusterHealth(target string) (result interface{
		ActivePrimaryShards() int     // get active primary shards number in cluster health result
		ActiveShards() int            // get active shards number in cluster health result
		UnassignedShards() int        // get unassigned shards number in cluster health result
		ActiveShardsPercent() float64 // get active shards percent in cluster health result
	}, err error)

	// GetIndicesWithRegexp return indices list with regexp pattern
	GetIndicesWithRegexp(pattern string) (indices []string, err error)

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
