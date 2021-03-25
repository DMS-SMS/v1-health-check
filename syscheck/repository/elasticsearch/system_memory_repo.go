// Create file in v.1.0.0
// syscheck_memory_repo.go is file that define repository implement about memory using elasticsearch
// memory repository struct embed esRepositoryRequiredComponent struct in ./syscheck.go file

package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
	"log"

	"github.com/DMS-SMS/v1-health-check/domain"
)

// esMemoryCheckHistoryRepository is to handle MemoryCheckHistory model using elasticsearch as data store
type esMemoryCheckHistoryRepository struct {
	// myCfg is used for get memory check history repository config about elasticsearch
	myCfg esMemoryCheckHistoryRepoConfig

	// esCli is elasticsearch client connection injected from the outside package
	esCli *elasticsearch.Client

	// bodyWriter is implementation of reqBodyWriter interface to write []byte for request body
	bodyWriter reqBodyWriter
}

// esMemoryCheckHistoryRepoConfig is the config for memory check history repository using elasticsearch
type esMemoryCheckHistoryRepoConfig interface {
	// get common method from embedding esRepositoryComponentConfig
	esRepositoryComponentConfig
}

// NewESMemoryCheckHistoryRepository return new object that implement MemoryCheckHistoryRepository interface
func NewESMemoryCheckHistoryRepository(cfg esMemoryCheckHistoryRepoConfig, cli *elasticsearch.Client, w reqBodyWriter) domain.MemoryCheckHistoryRepository {
	repo := &esMemoryCheckHistoryRepository{
		myCfg:      cfg,
		esCli:      cli,
		bodyWriter: w,
	}

	if err := repo.Migrate(); err != nil {
		log.Fatal(errors.Wrap(err, "could not migrate repository").Error())
	}

	return repo
}
