// Create package in v.1.0.0
// elasticsearch package is for implementations of syscheck domain repository using elasticsearch
// In practice, repository struct declaration and implementation is in this syscheck.go file

// syscheck.go is file that define structure to embed from another structures.
// It also defines variables or constants used jointly in the package as private.

package elasticsearch

import "io"

// constant for default value of esRepositoryComponent field
const (
	// default primary shard number of system check  index in elasticsearch
	defaultIndexShardNum = 1

	// default replica shard number of system check index in elasticsearch
	defaultIndexReplicaNum = 1

	// default name of system check index in elasticsearch
	defaultIndexName = "sms-system-check"
)

// esRepositoryComponent contains the least information that elasticsearch repository should have in syscheck
// If you change field name or type, you should check and change FieldSetter about that changed field in field_setter.go
type esRepositoryComponent struct {
	// indexName represent name of index including esDiskCheckHistory document
	IndexName string

	// indexShardNum represent shard number of index including esDiskCheckHistory document
	IndexShardNum int

	// indexReplicaNum represent replica number of index to replace index when node become unable
	IndexReplicaNum int
}

// FieldSetter is custom function type to set field value of interface with reflect package
type FieldSetter func(interface{})

// reqBodyWriter is private interface to use as writing []byte for request body
type reqBodyWriter interface {
	io.Writer
	io.WriterTo
}
