// Create file in v.1.0.0
// syscheck_cpu_ucase.go is file that define usecase implementation about syscheck cpu domain
// cpu check usecase struct embed systemCheckUsecaseComponent struct in ./syscheck.go file

package usecase

import (
	"context"
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

	// slackChat is used for agent slack API about chatting
	slackChatAgency slackChatAgency

	// cpuSysAgency is used as agency about cpu system command
	cpuSysAgency cpuSysAgency

	// status represent current process status of cpu health check
	status cpuCheckStatus

	// mutex help to prevent race condition when set status field value
	mutex sync.Mutex
}

// cpuCheckUsecaseConfig is the config getter interface for cpu check usecase
type cpuCheckUsecaseConfig interface {
	// get common config method from embedding systemCheckUsecaseComponentConfig
	systemCheckUsecaseComponentConfig

	// CPUWarningUsage method returns float64 represent cpu warning usage
	CPUWarningUsage() float64

	// CPUMaximumUsage method returns float64 represent cpu maximum usage
	CPUMaximumUsage() float64
}

// cpuSysAgency is agency that agent various command about cpu system
type cpuSysAgency interface {
	// CalculateContainersCPUUsage calculate container cpu usage & return result interface implementation
	CalculateContainersCPUUsage() (result interface{
		// TotalCPUUsage return total cpu usage in docker containers
		TotalCPUUsage() float64

		// MostConsumerExceptFor return container consume the most CPU except container names received from param
		MostConsumerExceptFor([]string) (id, name string, usage float64)
	}, err error)
}

// NewCPUCheckUsecase function return cpuCheckUsecase ptr instance after initializing
func NewCPUCheckUsecase(
	cfg cpuCheckUsecaseConfig,
	chr domain.CPUCheckHistoryRepository,
	sca slackChatAgency,
	csa cpuSysAgency,
) domain.CPUCheckUseCase {
	return &cpuCheckUsecase{
		// initialize field with parameter received from caller
		myCfg:           cfg,
		historyRepo:     chr,
		slackChatAgency: sca,
		cpuSysAgency:    csa,

		// initialize field with default value
		status: cpuStatusHealthy,
		mutex:  sync.Mutex{},
	}
}

// CheckCPU check cpu health with checkCPU method & store check history in repository
// Implement CheckCPU method of domain.CPUCheckUseCase interface
func (cu *cpuCheckUsecase) CheckCPU(ctx context.Context) error {
	history := cu.checkCPU(ctx)

	if b, err := cu.historyRepo.Store(history); err != nil {
		return errors.Wrapf(err, "failed to store cpu check history, response: %s", string(b))
	}

	return nil
}
