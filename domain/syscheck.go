// Create package in v.1.0.0
// Same as entities, struct and method in domain package will used in all layer.
// syscheck.go is file that define model as struct and abstract method of model as interface.
// Also, it declare usecase interface used as business layer.

// syscheck domain is managing the state of the system (CPU, memory, disk, etc.) periodically

// All model struct and interface is about 'system check' domain
// Most importantly, it only defines and does not implement interfaces.

package domain

import (
	"log"
	"os"
	"strings"
	"time"
)

// systemCheckHistoryComponent is basic model using by embedded in every model struct about check history
type systemCheckHistoryComponent struct {
	// private field in below, these fields have fixed value so always set in FillPrivateComponent method
	// Agent specifies name of service that created this model
	agent string

	// version specifies health checker version when this model was created
	version string

	// Timestamp specifies the time when this model was created.
	timestamp time.Time

	// Domain specifies domain about right this package, syscheck
	domain string

	// _type specifies detail service type in system check domain (Ex, DiskCheck, CPUCheck)
	_type string

	// ---

	// public field in below, these fields don't have fixed value so set in another package from custom user
	// UUID specifies Universally unique identifier in each of system check process
	UUID string

	// ProcessLevel specifies about how level to handle system check process.
	ProcessLevel processLevel

	// Message specifies additional description of result about system check process.
	Message string

	// Error specifies error message if health check's been handled abnormally.
	Error error

	// ---

	// field in below is about alarm result and is private so call SetAlarmResult method to set this field value
	// Alerted specifies if alert result or status in while handling system check process.
	alerted bool

	// alarmText specifies alarm text sent in system check process.
	alarmText string

	// alarmTime specifies time when this system check sent alarm.
	alarmTime time.Time

	// alarmErr specifies Error occurred when sending alarm.
	alarmErr error
}

// systemCheckHistoryRepositoryComponent is basic interface using by embedded in every repository about check history
type systemCheckHistoryRepositoryComponent interface {
	// Migrate method build environment for storage in stores such as Mysql or Elasticsearch, etc.
	Migrate() error
}

// FillComponent fill field of systemCheckHistoryComponent if is empty
func (sch *systemCheckHistoryComponent) FillPrivateComponent() {
	sch.version = version
	sch.agent = "sms-health-check"
	sch.domain = "syscheck"
	sch._type = "None"
	sch.timestamp = time.Now()
}

// DottedMapWithPrefix convert systemCheckHistoryComponent to dotted map and return that
// all key value of Map start with prefix received from parameter
func (sch *systemCheckHistoryComponent) DottedMapWithPrefix(prefix string) (m map[string]interface{}) {
	if prefix != "" {
		prefix += "."
	}

	m = map[string]interface{}{}

	// setting private field value in dotted map
	m[prefix + "version"] = sch.version
	m[prefix + "agent"] = sch.agent
	m[prefix + "@timestamp"] = sch.timestamp
	m[prefix + "domain"] = sch.domain
	m[prefix + "type"] = sch._type

	// setting public field value in dotted map
	m[prefix + "uuid"] = sch.UUID
	m[prefix + "process_level"] = sch.ProcessLevel.String()
	m[prefix + "message"] = sch.Message
	if sch.Error == nil {
		m[prefix + "error"] = nil
	} else {
		m[prefix + "error"] = sch.Error.Error()
	}

	// setting alarm result field value in dotted map
	m[prefix + "alerted"] = sch.alerted
	m[prefix + "alarm_text"] = sch.alarmText
	m[prefix + "alarm_time"] = sch.alarmTime
	m[prefix + "alarm_error"] = sch.alarmErr

	return
}

// SetAlarmResult set field value about alarm result with parameter
func (sch *systemCheckHistoryComponent) SetAlarmResult(t time.Time, text string, err error) {
	sch.alerted = true
	sch.alarmTime = t
	sch.alarmText = text
	sch.alarmErr = err
}

// processLevel is string custom type used for representing status check process level
type processLevel []string

// Set method overwrite processLevel slice to level received from parameter
func (pl *processLevel) Set(level string) {
	*pl = processLevel{level}
}

// Append method append processLevel slice with level received from parameter
func (pl *processLevel) Append(level string) {
	*pl = append(*pl, level)
}

// String method return string which join processLevel slice to string with " | "
func (pl *processLevel) String() string {
	return strings.Join(*pl, " | ")
}

// get information from system environment variable
var version string
func init() {
	if version = os.Getenv("VERSION"); version == "" {
		log.Fatal("please set VERSION in environment variable")
	}
}
