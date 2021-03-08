// Create package in v.1.0.0
// json package is a collection of convenience objects, specific interface implementations related to json.
// map_writer.go is file that declare customized objects called MapWriter to convert map-type variable to json format.
// Different things of mapWriter with json.Marshal() is that separate json value step with dot in map key.

package json

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"sort"
	"strings"
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

// Write method write bytes which is marshalling map[string]interface{} to field according to special rule.
// special rule is that separate json value step with dot in map key.
func (mw *mapWriter) Write(b []byte) (n int, err error) {
	// be unable to call Write method again before calling WriteTo method.
	mw.mu.Lock()

	var m map[string]interface{}
	if err = json.Unmarshal(b, &m); err != nil {
		err = errors.Wrap(err, "failed to json.Unmarshal bytes to map[string]interface{} type")
		mw.mu.Unlock()
		return
	}

	km := map[string]bool{}
	buf := map[string]interface{}{}
	ctx := map[string]interface{}{}

	// add front of key except last part with separated by dot in km(map[string]bool{}) and ctx(map[string]interface{}{})
	for k, v := range m {
		keys := strings.Split(k, ".")
		tail := keys[len(keys)-1]
		front := strings.Join(keys[:len(keys)-1], ".")

		// if keys length is 1, save that key in buf directly. And don't add key in km & ctx
		if len(keys) == 1 {
			buf[k] = v
			continue
		}

		km[front] = true
		if _, ok := ctx[front]; !ok {
			ctx[front] = map[string]interface{}{tail: v}
		} else {
			ctx[front].(map[string]interface{})[tail] = v
		}
	}

	// change map witch store front of key to string array in alignment with sort.Strings
	ks := make([]string, len(km))
	var i int
	for k := range km {
		ks[i] = k
		i++
	}
	sort.Strings(ks)

	// create derived maps with parsing key in ctx and put value in last map using agent
	for _, k := range ks {
		v := ctx[k]
		keys := strings.Split(k, ".")
		agent := &map[string]interface{}{}
		for i, k := range keys {
			if _, ok := buf[k]; i == 0 && ok {
				(*agent)[k] = buf[k]
			}

			if _, ok := (*agent)[k]; ok {
				// return error if cannot assert (*agent)[k] to *map[string]interface{} type
				// asserting error means that already last value with that key is exists.
				if agent, ok = (*agent)[k].(*map[string]interface{}); !ok {
					err = errors.Errorf("invalid key for json format, key: %s", strings.Join(keys, "."))
					mw.mu.Unlock()
					return
				}
			} else {
				(*agent)[k] = &map[string]interface{}{}
				agent = (*agent)[k].(*map[string]interface{})
				if i == 0 {
					buf[k] = agent
				}
			}

			if (i + 1) == len(keys) {
				*agent = v.(map[string]interface{})
			}
		}
	}

	if b, err = json.Marshal(buf); err != nil {
		err = errors.Wrap(err, "failed to json.Marshal map[string]interface{} buffer")
		mw.mu.Unlock()
		return
	}

	mw.buf = b
	mw.written = true
	n = len(b)
	return
}

// WriteTo method write bytes buf created by Write method to Writer received from parameter.
// Cannot use WriteTo method before calling Write method.
func (mw *mapWriter) WriteTo(w io.Writer) (n int64, err error) {
	// return error if Write method is not called before calling WriteTo method.
	if !mw.written {
		err = errors.New("you should call Write method before calling WriteTo method")
		return
	}

	// be able to call Write method after this method finished
	defer mw.mu.Unlock()

	mw.written = false
	if _, err = w.Write(mw.buf); err != nil {
		err = errors.Wrap(err, "failed to write to Writer received from parameter")
		return
	}
	n = int64(len(mw.buf))
	return
}
