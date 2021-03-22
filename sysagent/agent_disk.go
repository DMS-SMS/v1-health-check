// Create file in v.1.0.0
// agent_disk.go is file that define method of sysAgent that agent command about disk
// For example in disk command, there are get remaining capacity, prune disk, etc ...

package sysagent

import (
	"github.com/inhies/go-bytesize"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
	"os"
)

// GetRemainDiskCapacity return remain disk capacity expressed in bytesize package
func (sa *sysAgent) GetRemainDiskCapacity() (size bytesize.ByteSize, err error) {
	var stat unix.Statfs_t

	wd, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "failed to call os.Getwd")
		return
	}

	if err = unix.Statfs(wd, &stat); err != nil {
		err = errors.Wrap(err, "failed to call unix.Statfs")
		return
	}

	// Available blocks * size per block = available space in bytes
	size = bytesize.New(float64(stat.Bavail * uint64(stat.Bsize)))
	return
}
