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
	// Agent specifies name of service that created this model
	Agent string

	// version specifies health checker version when this model was created
	Version string

	// Timestamp specifies the time when this model was created.
	Timestamp time.Time

	// Domain specifies domain about right this package, syscheck
	Domain string

	// Type specifies detail service type in system check domain (Ex, DiskCheck, CPUCheck)
	Type string

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
func (sch *systemCheckHistoryComponent) FillComponent() {
	if sch.Version == "" {
		sch.Version = version
	}

	if sch.Agent == "" {
		sch.Agent = "sms-health-check"
	}

	now := time.Now()
	if now.Location().String() == time.UTC.String() {
		now = now.Add(time.Hour * 9)
	}
	sch.Timestamp = now

	if sch.Domain == "" {
		sch.Domain = "syscheck"
}
}

// MapWithPrefixKey convert systemCheckHistoryComponent to dotted map and return that
// all key value of Map start with prefix received from parameter
func (sch *systemCheckHistoryComponent) MapWithPrefixKey(prefix string) (m map[string]interface{}) {
	if prefix != "" {
		prefix += "."
	}

	m = map[string]interface{}{}
	m[prefix + "version"] = sch.Version
	m[prefix + "agent"] = sch.Agent
	m[prefix + "@timestamp"] = sch.Timestamp
	m[prefix + "type"] = sch.Type
	m[prefix + "process_level"] = sch.ProcessLevel
	m[prefix + "alerted"] = sch.Alerted
	m[prefix + "alarm_content"] = sch.AlarmContent

	return
}
