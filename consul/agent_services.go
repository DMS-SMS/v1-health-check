// Create file in v.1.0.0
// agent_services.go is file that define method of consulAgent that agent command about services
// For example in consul command, there are get services, deregister service, etc ...

package consul

import (
	"github.com/pkg/errors"
)

// GetAllServices method get all services in consul & return services interface implement
func (ca *consulAgent) GetAllServices() (interface {
	IDsOf(srv string) (ids []string) // idsOf return id list of instances which are of received srv
	All() (srvs map[string][]string) // All return all id list of all services as map
}, error) {
	srvs, err := ca.cslCli.Agent().Services()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all services in consul")
	}

	srvM := services{}
	for _, srv := range srvs {
		if _, ok := srvM[srv.Service]; !ok {
			srvM[srv.Service] = []string{}
		}
		srvM[srv.Service] = append(srvM[srv.Service], srv.ID)
	}
	return srvM, nil
}

// DeregisterInstance method deregister instance in consul with received id
func (ca *consulAgent) DeregisterInstance(id string) (err error) {
	return errors.Wrap(ca.cslCli.Agent().ServiceDeregister(id), "failed to deregister consul service")
}

// services is map binding type having id list per services, and implement GetAllServices return type interface
type services map[string][]string

func (s services) IDsOf(srv string) (ids []string) { ids, _ = s[srv]; return } // IDsOf return id list of service
func (s services) All() (srvs map[string][]string) { return s }                // All return all services id list
