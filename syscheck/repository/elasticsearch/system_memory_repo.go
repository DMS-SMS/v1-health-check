// Create file in v.1.0.0
// syscheck_memory_repo.go is file that define repository implement about memory using elasticsearch
// disk repository struct embed esRepositoryRequiredComponent struct in ./syscheck.go file

package elasticsearch

import "github.com/elastic/go-elasticsearch/v7"

// esMemoryCheckHistoryRepository is to handle MemoryCheckHistory model using elasticsearch as data store
type esMemoryCheckHistoryRepository struct {
	// myCfg is used for get disk check history repository config about elasticsearch
	myCfg esMemoryCheckHistoryRepoConfig

	// esCli is elasticsearch client connection injected from the outside package
	esCli *elasticsearch.Client

	// bodyWriter is implementation of reqBodyWriter interface to write []byte for request body
	bodyWriter reqBodyWriter
}
