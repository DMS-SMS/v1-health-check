// Create file in v.1.0.0
// syscheck_memory_ucase.go is file that define usecase implementation about syscheck memory domain
// memory check usecase struct embed systemCheckUsecaseComponent struct in ./syscheck.go file

package usecase

// memoryCheckStatus is type to int constant represent current memory check process status
type memoryCheckStatus int
const (
	memoryStatusHealthy    memoryCheckStatus = iota // represent memory check status is healthy
	memoryStatusWarning                             // represent memory check status is warning now
	memoryStatusRecovering                          // represent it's recovering memory status now
	memoryStatusUnhealthy                           // represent memory check status is unhealthy
)
