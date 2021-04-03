// Create file in v.1.0.0
// srvcheck_swarmpit_ucase.go is file that define usecase implementation about swarmpit check in srvcheck domain
// usecase layer depend on repository layer and is depended to delivery layer

package usecase

// swarmpitCheckStatus is type to int constant represent current swarmpit check process status
type swarmpitCheckStatus int
const (
	swarmpitStatusHealthy    swarmpitCheckStatus = iota // represent swarmpit check status is healthy
	swarmpitStatusRecovering                            // represent it's recovering swarmpit status now
	swarmpitStatusUnhealthy                             // represent swarmpit check status is unhealthy
)
