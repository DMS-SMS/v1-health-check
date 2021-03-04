// Create package in v.1.0.0
// elasticsearch package is for implementations of syscheck domain repository using elasticsearch
// In practice, repository struct declaration and implementation is in this syscheck.go file

package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v7"
)

const (
	// default primary shard number of system check  index in elasticsearch
	defaultIndexShardNum = 1

	// default replica shard number of system check index in elasticsearch
	defaultIndexReplicaNum = 1

	// default name of system check index in elasticsearch
	defaultIndexName = "sms-system-check"
)

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
