// Create file in v.1.0.0
// agent_memory.go is file that define method of sysAgent that agent command about memory
// For example in memory command, there are get total memory usage, calculate container memory usage, etc ...

package system

import (
	"github.com/inhies/go-bytesize"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/pkg/errors"
)

// GetTotalSystemMemoryUsage return total memory usage as bytesize in system
func (sa *sysAgent) GetTotalSystemMemoryUsage() (size bytesize.ByteSize, err error) {
	stats, err := memory.Get()
	if err != nil {
		err = errors.Wrap(err, "failed to get memory stats")
		return
	}

	size = bytesize.ByteSize(stats.Used)
	return
}
