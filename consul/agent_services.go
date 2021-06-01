// Create file in v.1.0.0
// agent_services.go is file that define method of consulAgent that agent command about services
// For example in consul command, there are get services, deregister service, etc ...

package consul

import (
	"fmt"
	"github.com/pkg/errors"
)

// GetServices method get services in consul & return services interface implement
func (ca *consulAgent) GetServices(srv string) (interface {
	HasNext() bool           // HasNext method return if srvIter has next element
	Next() (id, addr string) // Next method return next service id, address
}, error) {
	srvs, err := ca.cslCli.Agent().ServicesWithFilter(fmt.Sprintf("Service==%s", srv))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all services in consul")
	}

	iter := &srvIter{
		srv: []struct{ id, addr string }{},
		x:   0,
	}

	for _, srv := range srvs {
		iter.srv = append(iter.srv, struct {
			id, addr string
		}{
			id:   srv.ID,
			addr: fmt.Sprintf("%s:%d", srv.Address, srv.Port),
		})
	}

	return iter, nil
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
func (si *srvIter) HasNext() bool {
	return si.x < len(si.srv)
}

// Next method return next service id, address
func (si *srvIter) Next() (id, addr string) {
	id = si.srv[si.x].id
	addr = si.srv[si.x].addr
	si.x++
	return
}
