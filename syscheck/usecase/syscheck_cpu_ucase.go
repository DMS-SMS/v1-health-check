// Create file in v.1.0.0
// syscheck_cpu_ucase.go is file that define usecase implementation about syscheck cpu domain
// cpu check usecase struct embed systemCheckUsecaseComponent struct in ./syscheck.go file

package usecase

import (
	"github.com/docker/docker/client"
	"sync"

	"github.com/DMS-SMS/v1-health-check/domain"
)

// cpuCheckStatus is type to int constant represent current cpu check process status
type cpuCheckStatus int
const (
	cpuStatusHealthy    cpuCheckStatus = iota // represent cpu check status is healthy
	cpuStatusRecovering                       // represent it's recovering cpu status now
	cpuStatusUnhealthy                        // represent cpu check status is unhealthy
)

// cpuCheckUsecase implement CPUCheckUsecase interface in domain and used in delivery layer
type cpuCheckUsecase struct {
	// myCfg is used for getting cpu check usecase config
	myCfg cpuCheckUsecaseConfig

	// historyRepo is used for store cpu check history and injected from outside
	historyRepo domain.CPUCheckHistoryRepository

	// dockerCli is docker client to call docker agent API
	dockerCli *client.Client

	// slackChat is used for agent slack API about chatting
	slackChatAgency slackChatAgency

	// status represent current process status of cpu health check
	status cpuCheckStatus

	// mutex help to prevent race condition when set status field value
	mutex sync.Mutex
}

// cpuCheckUsecaseConfig is the config getter interface for cpu check usecase
type cpuCheckUsecaseConfig interface {
	// get common config method from embedding systemCheckUsecaseComponentConfig
	systemCheckUsecaseComponentConfig

	// CPUUsageLimit method returns float32 represent cpu usage limit value
	CPUUsageLimit() float32
}