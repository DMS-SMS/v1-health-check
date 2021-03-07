// Create package in v.1.0.0
// json package is a collection of convenience objects, specific interface implementations related to json.
// map_writer.go is file that declare customized objects called MapWriter to convert map-type variable to json format.
// Different things of mapWriter with json.Marshal() is that separate json value step with dot in map key .

package json

import (
	"github.com/pkg/errors"
	"io"
	"sync"
)

// mapWriter is struct to write []byte which is map type marshaled to another []byte for json
// mapWriter implement io.Writer & io.WriterTo interface
type mapWriter struct {
	// buf is field that stores created bytes[] through Write method
	buf []byte

	// written is used for checking if Write method called before calling WriteTo method
	written bool

	// mu prevents other gorouting from entering after Write method call until WriteTo method call
	mu sync.Mutex
}

// MapWriter return new pointer instance of mapWriter struct
func MapWriter() *mapWriter {
	return &mapWriter{
		buf:     []byte{},
		written: false,
		mu:      sync.Mutex{},
	}
}

// WriteTo method write bytes buf created by Write method to Writer received from parameter.
// Cannot use WriteTo method before calling Write method.
func (mw *mapWriter) WriteTo(w io.Writer) (n int64, err error) {
	// return error if Write method is not called before calling WriteTo method.
	if !mw.written {
		err = errors.New("you should call Write method before calling WriteTo method")
		return
	}

	// able to call Write method after this method finished
	defer mw.mu.Unlock()

	mw.written = false
	if _, err = w.Write(mw.buf); err != nil {
		err = errors.Wrap(err, "failed to write to Writer received from parameter")
		return
	}
	n = int64(len(mw.buf))
	return
}
