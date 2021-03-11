// Create package in v.1.0.0
// elasticsearch package is for implementations of syscheck domain repository using elasticsearch
// In practice, repository struct declaration and implementation is in this syscheck.go file

// syscheck.go is file that define structure to embed from another structures.
// It also defines interface or function used jointly in the package as private.

package elasticsearch

import "io"

// esRepositoryComponentConfig is interface contains method to return config value that elasticsearch repository should have
// It can be externally set as Config object that implements that interface.
type esRepositoryComponentConfig interface {
	// IndexName method returns the index name of elasticsearch about syscheck
	IndexName() string

	// IndexShardNum method returns the number of index shard in elasticsearch about syscheck
	IndexShardNum() int

	// IndexReplicaNum method returns the number of index replica in elasticsearch about syscheck
	IndexReplicaNum() int
}

// reqBodyWriter is private interface to use as writing []byte for request body
type reqBodyWriter interface {
	io.Writer
	io.WriterTo
}
