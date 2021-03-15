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
	"time"
)

// get information from system environment variable
var version string
func init() {
	if version = os.Getenv("VERSION"); version == "" {
		log.Fatal("please set VERSION in environment variable")
	}
}

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


	// public field in below, these fields don't have fixed value so set in another package from custom user
	// ProcessLevel specifies about how level to handle system check process.
	ProcessLevel string

	// Alerted specifies if alert result or status in while handling system check process.
	Alerted bool

	// AlarmContent specifies content about alarm sent in system check process.
	AlarmContent string
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

	now := time.Now()
	if now.Location().String() == time.UTC.String() {
		now = now.Add(time.Hour * 9)
	}
	sch.timestamp = now
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

	// setting public field value in dotted map
	m[prefix + "type"] = sch.Type
	m[prefix + "process_level"] = sch.ProcessLevel
	m[prefix + "alerted"] = sch.Alerted
	m[prefix + "alarm_content"] = sch.AlarmContent

	return
}
