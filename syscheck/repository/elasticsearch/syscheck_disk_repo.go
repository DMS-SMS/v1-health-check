// Create file in v.1.0.0
// syscheck_disk_repo.go is file that define repository implement about disk using elasticsearch
// disk repository struct embed esRepositoryRequiredComponent struct in ./syscheck.go file

package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/pkg/errors"
)

// esDiskCheckHistoryRepository is to handle DiskCheckHistory model using elasticsearch as data store
type esDiskCheckHistoryRepository struct {
	// get common field from embedding esRepositoryComponent
	esRepositoryComponent

	// esCli is elasticsearch client connection injected from the outside package
	esCli *elasticsearch.Client

	// bodyWriter is implementation of reqBodyWriter interface to write []byte for request body
	bodyWriter reqBodyWriter
}

// NewESDiskCheckHistoryRepository return new object that implement DiskCheckHistory.Repository interface
func NewESDiskCheckHistoryRepository(cli *elasticsearch.Client, w reqBodyWriter, setters ...FieldSetter) (domain.DiskCheckHistoryRepository, error) {
	repo := &esDiskCheckHistoryRepository{
		esCli:      cli,
		bodyWriter: w,
	}

	repo.IndexName = defaultIndexName
	repo.IndexShardNum = defaultIndexShardNum
	repo.IndexReplicaNum = defaultIndexReplicaNum

	// set repository field by running FieldSetter type functions received from caller
	for _, s := range setters {
		s(repo)
	}

	if err := repo.Migrate(); err != nil {
		return nil, errors.Wrap(err, "could not migrate repository")
	}

	return repo, nil
}
