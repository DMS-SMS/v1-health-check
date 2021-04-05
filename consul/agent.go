// Create package in v.1.0.0
// consul package define struct which is implement various interface about consul agency using in each of domain
// there are kind of method in consul agency such as get services, deregister service, etc ...

// in agent.go file, define struct type of consul agent & initializer that are not method.
// Also if exist, custom type or variable used in common in each of method will declared in this file.

package consul

import (
	"github.com/hashicorp/consul/api"
)

// consulAgent is struct that agent various command about consul including get services, deregister service, etc ...
type consulAgent struct {
	// cslCli is client connection about consul & can access consul API with this client
	cslCli *api.Client
}

// NewAgent return new instance of consulAgent pointer type initialized with parameter
func NewAgent(cc *api.Client) *consulAgent {
	return &consulAgent{
		cslCli: cc,
	}
}
