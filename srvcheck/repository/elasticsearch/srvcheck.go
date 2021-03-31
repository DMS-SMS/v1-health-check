// Create package in v.1.0.0
// elasticsearch package is for implementations of srvcheck domain repository using elasticsearch
// In practice, repository struct declaration and implementation occur in this package

// srvcheck.go is file that define structure to embed from another structures.
// It also defines interface or function used jointly in the package as private.

package elasticsearch

// esRepositoryComponentConfig is interface contains method to return config value that elasticsearch repository should have
// It can be externally set as Config object that implements that interface.
type esRepositoryComponentConfig interface {
	// IndexName method returns the index name of elasticsearch about srvcheck
	IndexName() string

	// IndexShardNum method returns the number of index shard in elasticsearch about srvcheck
	IndexShardNum() int

	// IndexReplicaNum method returns the number of index replica in elasticsearch about srvcheck
	IndexReplicaNum() int
}
