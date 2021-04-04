// Create file in v.1.0.0
// srvcheck_consul_repo.go is file that define implement consul history repository using elasticsearch
// this elasticsearch repository struct embed esRepositoryRequiredComponent struct in ./srvcheck.go file

package elasticsearch

import "github.com/elastic/go-elasticsearch/v7"

// esConsulCheckHistoryRepository is to handle ConsulCheckHistoryRepository model using elasticsearch as data store
type esConsulCheckHistoryRepository struct {
	// myCfg is used for get consul history repository config about elasticsearch
	myCfg esConsulCheckHistoryRepoConfig

	// esCli is elasticsearch client connection injected from the outside package
	esCli *elasticsearch.Client

	// reqBodyWriter is implementation of reqBodyWriter interface to write []byte for request body
	reqBodyWriter reqBodyWriter
}
