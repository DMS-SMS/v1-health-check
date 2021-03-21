// Create package in v.1.0.0
// measure package define measurer struct about system, msg(문자 서비스) usage, etc ...
// system_measurer.go define system measurer about cpu or memory usage, disk remain capacity, etc ...

package measure

import (
	"github.com/inhies/go-bytesize"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
	"os"
)

// systemMeasurer is struct that measure value about system
type systemMeasurer struct {}

// SystemMeasurer function return systemMeasurer ptr instance with initializing
func SystemMeasurer() *systemMeasurer {
	return &systemMeasurer{}
}

// RemainDiskCapacity return remain disk capacity expressed in bytesize package
func (sm *systemMeasurer) RemainDiskCapacity() (size bytesize.ByteSize, err error) {
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
