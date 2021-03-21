// Create file in v.1.0.0
// syscheck_cpu_ucase.go is file that define usecase implementation about syscheck cpu domain
// cpu check usecase struct embed systemCheckUsecaseComponent struct in ./syscheck.go file

package usecase

// cpuCheckStatus is type to int constant represent current cpu check process status
type cpuCheckStatus int
const (
	cpuStatusHealthy    cpuCheckStatus = iota // represent cpu check status is healthy
	cpuStatusRecovering                       // represent it's recovering cpu status now
	cpuStatusUnhealthy                        // represent cpu check status is unhealthy
)

// cpuCheckUsecaseConfig is the config getter interface for cpu check usecase
type cpuCheckUsecaseConfig interface {
	// get common config method from embedding systemCheckUsecaseComponentConfig
	systemCheckUsecaseComponentConfig

	// CPUUsageLimit method returns float32 represent cpu usage limit value
	CPUUsageLimit() float32
}
