// Create file in v.1.0.0
// srvcheck_consul_ucase.go is file that define usecase implementation about consul check in srvcheck domain
// usecase layer depend on repository layer and is depended to delivery layer

package usecase

// consulCheckStatus is type to int constant represent current consul check process status
type consulCheckStatus int
const (
	consulStatusHealthy    consulCheckStatus = iota // represent consul check status is healthy
	consulStatusRecovering                          // represent it's recovering consul status now
	consulStatusUnhealthy                           // represent consul check status is unhealthy
)
