// Create file in v.1.0.0
// agent_container.go is file that define method of dockerAgent that agent command about docker container
// For example in disk command, there are remove or restart container, etc ...

package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
)

// RemoveContainer remove container with id & option (auto created from docker swarm if exists)
func (da *dockerAgent) RemoveContainer(containerID string, options types.ContainerRemoveOptions) error {
	var (
		ctx = context.Background()
	)

	return errors.Wrap(da.dkrCli.ContainerRemove(ctx, containerID, options), "failed to call ContainerRemove")
}
