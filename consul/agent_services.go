// Create file in v.1.0.0
// agent_services.go is file that define method of consulAgent that agent command about services
// For example in consul command, there are get services, deregister service, etc ...

package consul

// services is map binding type having id list per services, and implement GetAllServices return type interface
type services map[string][]string

// idsOf return id list of instances which are of received srv
func (s services) IDsOf(srv string) (ids []string) { ids, _ = s[srv]; return }
