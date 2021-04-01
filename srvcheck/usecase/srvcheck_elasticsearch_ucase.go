// Create file in v.1.0.0
// srvcheck_elasticsearch_ucase.go is file that define usecase implementation about elasticsearch check in srvcheck domain
// elasticsearch check usecase struct embed serviceCheckUsecaseComponent struct in ./srvcheck.go file

package usecase

import (
	"github.com/DMS-SMS/v1-health-check/domain"
	"sync"
	"time"
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

	// status represent current process status of elasticsearch health check
	status elasticsearchCheckStatus

	// mutex help to prevent race condition when set status field value
	mutex sync.Mutex
}

// elasticsearchCheckUsecaseConfig is the config getter interface for elasticsearch check usecase
type elasticsearchCheckUsecaseConfig interface {
	// get common config method from embedding serviceCheckUsecaseComponentConfig
	serviceCheckUsecaseComponentConfig

	// TargetClusterName method returns string represent target cluster name
	TargetClusterName() string

	// MaximumShardsNumber method returns int represent maximum shards number
	MaximumShardsNumber() int

	// JaegerIndexRegexp method returns string represent jaeger index regexp
	JaegerIndexRegexp() string

	// JaegerIndexLifeCycle method returns duration represent jaeger index life cycle
	JaegerIndexLifeCycle() time.Duration
}
