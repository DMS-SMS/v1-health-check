// Create file in v.1.0.0
// srvcheck_elasticsearch_repo.go is file that define implement elasticsearch history repository using elasticsearch
// this elasticsearch repository struct embed esRepositoryRequiredComponent struct in ./srvcheck.go file

package elasticsearch

import (
	"github.com/DMS-SMS/v1-health-check/domain"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
	"log"
)

// esElasticsearchCheckHistoryRepository is to handle ElasticsearchCheckHistoryRepository model using elasticsearch as data store
type esElasticsearchCheckHistoryRepository struct {
	// myCfg is used for get cpu elasticsearch history repository config about elasticsearch
	myCfg esElasticsearchCheckHistoryRepoConfig

	// esCli is elasticsearch client connection injected from the outside package
	esCli *elasticsearch.Client

	// reqBodyWriter is implementation of reqBodyWriter interface to write []byte for request body
	reqBodyWriter reqBodyWriter
}

// esElasticsearchCheckHistoryRepoConfig is the config for elasticsearch check history repository using elasticsearch
type esElasticsearchCheckHistoryRepoConfig interface {
	// get common method from embedding esRepositoryComponentConfig
	esRepositoryComponentConfig
}

// NewESElasticsearchCheckHistoryRepository return new object that implement ElasticsearchCheckHistoryRepository interface
func NewESElasticsearchCheckHistoryRepository(
	cfg esElasticsearchCheckHistoryRepoConfig,
	cli *elasticsearch.Client,
	w reqBodyWriter,
) domain.ElasticsearchCheckHistoryRepository {
	repo := &esElasticsearchCheckHistoryRepository{
		myCfg:         cfg,
		esCli:         cli,
		reqBodyWriter: w,
	}

	if err := repo.Migrate(); err != nil {
		log.Fatal(errors.Wrap(err, "could not migrate repository").Error())
	}

	return repo
}