// Create file in v.1.0.0
// srvcheck_elasticsearch_ucase.go is file that define usecase implementation about elasticsearch check in srvcheck domain
// elasticsearch check usecase struct embed serviceCheckUsecaseComponent struct in ./srvcheck.go file

package usecase

// elasticsearchCheckStatus is type to int constant represent current elasticsearch check process status
type elasticsearchCheckStatus int
const (
	elasticsearchStatusHealthy    elasticsearchCheckStatus = iota // represent elasticsearch check status is healthy
	elasticsearchStatusWarning                                    // represent elasticsearch check status is warning now
	elasticsearchStatusRecovering                                 // represent it's recovering elasticsearch status now
	elasticsearchStatusUnhealthy                                  // represent elasticsearch check status is unhealthy
)
