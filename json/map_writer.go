// Create package in v.1.0.0
// json package is a collection of convenience objects, specific interface implementations related to json
// map_writer.go is file that declare customized objects called MapWriter to convert map-type variable to json format.

package json

import (
	"sync"
)

// mapWriter is struct to write []byte which is map type marshaled to another []byte for json
// mapWriter implement io.Writer & io.WriterTo interface
type mapWriter struct {
	// buf is field that stores created bytes[] through Write method
	buf []byte

	// mu prevents other gorouting from entering after Write method call until WriteTo method call
	mu sync.Mutex
}

// MapWriter return new pointer instance of mapWriter struct
func MapWriter() *mapWriter {
	return &mapWriter{
		buf: []byte{},
		mu:  sync.Mutex{},
	}
}
