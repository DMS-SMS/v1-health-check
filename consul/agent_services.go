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
type srvIter struct {
	srv []struct{ id, addr string } // srv contains array of struct having id, addr
	x   int                         // x represent current access index in iterator
}

// HasNext method return if srvIter has next element
func(si *srvIter) HasNext() bool {
	return si.x < len(si.srv)
}

// Next method return next service id, address
func(si *srvIter) Next() (id, addr string) {
	id = si.srv[si.x].id
	addr = si.srv[si.x].addr
	si.x++
	return
}
