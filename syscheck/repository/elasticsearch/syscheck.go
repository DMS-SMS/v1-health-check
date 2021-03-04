// Create package in v.1.0.0
// elasticsearch package is for implementations of syscheck domain repository using elasticsearch
// In practice, repository struct declaration and implementation is in this syscheck.go file

package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v7"
)

// esRepository contains the least information that elasticsearch repository should have in syscheck
type esRepository struct {
	// indexName represent name of index including esDiskCheckHistory document
	indexName string

	// indexShardNum represent shard number of index including esDiskCheckHistory document
	indexShardNum int

	// indexReplicaNum represent replica number of index to replace index when node become unable
	indexReplicaNum int
}

// esDiskCheckHistoryRepository is to handle DiskCheckHistory model using elasticsearch as data store
type esDiskCheckHistoryRepository struct {
	cli *elasticsearch.Client
}
