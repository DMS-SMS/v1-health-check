// Create file in v.1.0.0
// srvcheck_swarmpit_repo.go is file that define implement swarmpit history repository using elasticsearch
// this elasticsearch repository struct embed esRepositoryRequiredComponent struct in ./srvcheck.go file

package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
	"log"

	"github.com/DMS-SMS/v1-health-check/domain"
)

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

// NewESSwarmpitCheckHistoryRepository return new object that implement SwarmpitCheckHistoryRepository interface
func NewESSwarmpitCheckHistoryRepository(
	cfg esSwarmpitCheckHistoryRepoConfig,
	cli *elasticsearch.Client,
	w reqBodyWriter,
) domain.SwarmpitCheckHistoryRepository {
	repo := &esSwarmpitCheckHistoryRepository{
		myCfg:         cfg,
		esCli:         cli,
		reqBodyWriter: w,
	}

	if err := repo.Migrate(); err != nil {
		log.Fatal(errors.Wrap(err, "could not migrate repository").Error())
	}

	return repo
}
