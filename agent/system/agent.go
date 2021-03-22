// Create package in v.1.0.0
// sysagent package define system agency that agent various command about system
// For example about disk system, agent get remaining disk capacity, prune disk with docker, etc ...

// in agent.go file, define agent struct & initializer
// this agent implement various interface about system agency defined and using in usecase layer of syscheck domain

package system

import "github.com/docker/docker/client"

// sysAgent is struct that agent various command about system of disk, cpu, memory, etc ...
type sysAgent struct {
	// dockerCli is docker client to call docker agent API
	dockerCli *client.Client

	// containersCPUUsage is to keep cpu usage each of container get from GetTotalCPUUsage
	containersCPUUsage []struct {
		id, name string
		cpuUsage float64
	}
}

// New return new instance of sysAgent pointer type initialized with parameter
func New(dc *client.Client) *sysAgent {
	return &sysAgent{
		dockerCli: dc,
	}
}
