// Create package in v.1.0.0
// Same as entities, struct and method in domain package will used in all layer.
// srvcheck.go is file that define model as struct and abstract method of model as interface.
// Also, it declare usecase interface used as business layer.

// srvcheck domain is managing the state of the service (elasticsearch, swarm, consul, etc ...) periodically

// All model struct and interface is about service check domain
// Most importantly, it only defines and does not implement interfaces.

package domain

import (
	"strings"
	"time"
)

// serviceCheckHistoryComponent is basic model using by embedded in every model struct about service check history
type serviceCheckHistoryComponent struct {
	// private field in below, these fields have fixed value so always set in FillPrivateComponent method
	// Agent specifies name of service that created this model
	agent string

	// version specifies health checker version when this model was created
	version string

	// Timestamp specifies the time when this model was created.
	timestamp time.Time

	// Domain specifies domain about right this package, srvcheck
	domain string

	// _type specifies detail service type in service check domain (Ex, elasticsearch, swarm, consul, etc ...)
	_type string

	// ---

	// public field in below, these fields don't have fixed value so set in another package from custom user
	// UUID specifies Universally unique identifier in each of service check process
	UUID string

	// ProcessLevel specifies about how level to handle service check process.
	ProcessLevel srvcheckProcessLevel

	// Message specifies additional description of result about service check process.
	Message string

	// Error specifies error message if health check's been handled abnormally.
	Error error

	// ---

	// field in below is about alarm result and is private so call SetAlarmResult method to set this field value
	// Alerted specifies if alert result or status in while handling service check process.
	alerted bool

	// alarmText specifies alarm text sent in service check process.
	alarmText string

	// alarmTime specifies time when this service check sent alarm.
	alarmTime time.Time

	// alarmErr specifies Error occurred when sending alarm.
	alarmErr error
}

// FillComponent fill field of systemCheckHistoryComponent if is empty
func (sch *serviceCheckHistoryComponent) FillPrivateComponent() {
	sch.version = version
	sch.agent = "sms-health-check"
	sch.domain = "srvcheck"
	sch._type = "None"
	sch.timestamp = time.Now()
}

// srvcheckProcessLevel is string custom type used for representing service check process level
type srvcheckProcessLevel []string

// Set method overwrite srvcheckProcessLevel slice to level received from parameter
func (pl *srvcheckProcessLevel) Set(level string) {
	*pl = srvcheckProcessLevel{level}
}

// Append method append srvcheckProcessLevel slice with level received from parameter
func (pl *srvcheckProcessLevel) Append(level string) {
	*pl = append(*pl, level)
}

// String method return string which join srvcheckProcessLevel slice to string with " | "
func (pl *srvcheckProcessLevel) String() string {
	return strings.Join(*pl, " | ")
}
