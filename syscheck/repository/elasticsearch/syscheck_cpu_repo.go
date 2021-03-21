// Create file in v.1.0.0
// syscheck_cpu_repo.go is file that define implement cpu history repository using elasticsearch
// this cpu repository struct embed esRepositoryRequiredComponent struct in ./syscheck.go file

package elasticsearch

import "github.com/elastic/go-elasticsearch/v7"

// esCPUCheckHistoryRepository is to handle CPUCheckHistory model using elasticsearch as data store
type esCPUCheckHistoryRepository struct {
	// myCfg is used for get cpu check history repository config about elasticsearch
	myCfg esCPUCheckHistoryRepoConfig

	// esCli is elasticsearch client connection injected from the outside package
	esCli *elasticsearch.Client

	// bodyWriter is implementation of reqBodyWriter interface to write []byte for request body
	bodyWriter reqBodyWriter
}

// esCPUCheckHistoryRepoConfig is the config for cpu check history repository using elasticsearch
type esCPUCheckHistoryRepoConfig interface {
	// get common method from embedding esRepositoryComponentConfig
	esRepositoryComponentConfig
}
