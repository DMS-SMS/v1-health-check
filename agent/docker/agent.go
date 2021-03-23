// Create package in v.1.0.0
// docker package define struct which is implement various interface about docker agency using in each of domain
// there are kind of method in docker agency such as container, service, etc ...

// in agent.go file, define struct type of docker agent & initializer that are not method.
// Also if exist, custom type or variable used in common in each of method will declared in this file.

package docker

import "github.com/docker/docker/client"

// dockerAgent is struct that agent various command about docker including container, service, etc ...
type dockerAgent struct {
	dkrCli *client.Client
}

// New return new instance of dockerAgent pointer type initialized with parameter
func New(dc *client.Client) *dockerAgent {
	return &dockerAgent{
		dkrCli: dc,
	}
}
