// Create file in v.1.0.0
// srvcheck_swarmpit_repo.go is file that define implement swarmpit history repository using elasticsearch
// this elasticsearch repository struct embed esRepositoryRequiredComponent struct in ./srvcheck.go file

package elasticsearch

import "github.com/elastic/go-elasticsearch/v7"

// esSwarmpitCheckHistoryRepository is to handle SwarmpitCheckHistoryRepository model using elasticsearch as data store
type esSwarmpitCheckHistoryRepository struct {
	// myCfg is used for get swarmpit history repository config about elasticsearch
	myCfg esSwarmpitCheckHistoryRepoConfig

	// esCli is elasticsearch client connection injected from the outside package
	esCli *elasticsearch.Client

	// reqBodyWriter is implementation of reqBodyWriter interface to write []byte for request body
	reqBodyWriter reqBodyWriter
}

// esSwarmpitCheckHistoryRepoConfig is the config for swarmpit check history repository using elasticsearch
type esSwarmpitCheckHistoryRepoConfig interface {
	// get common method from embedding esRepositoryComponentConfig
	esRepositoryComponentConfig
}
