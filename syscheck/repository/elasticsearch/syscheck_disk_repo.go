// Create file in v.1.0.0
// syscheck_disk_repo.go is file that define repository implement about disk using elasticsearch
// disk repository struct embed esRepositoryRequiredComponent struct in ./syscheck.go file

package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/DMS-SMS/v1-health-check/domain"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
	"net/http"
	"time"
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

// Implement Migrate method of DiskCheckHistoryRepository interface
func (edr *esDiskCheckHistoryRepository) Migrate() error {
	resp, err := (esapi.IndicesExistsRequest{
		Index: []string{"gateway"},
	}).Do(context.Background(), edr.esCli)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to call IndicesExists, resp: %+v", resp))
	}

	if resp.StatusCode == http.StatusNotFound {
		if err := edr.createIndex(); err != nil {
			return errors.Wrap(err, "failed to create index")
		}
	}

	return nil
}

// createIndex method create index with name, share number in esDiskCheckHistoryRepository
func (edr *esDiskCheckHistoryRepository) createIndex() error {
	body := map[string]interface{}{}
	body["settings.number_of_shards"] = edr.IndexShardNum
	body["settings.number_of_replicas"] = edr.IndexReplicaNum

	b, _ := json.Marshal(body)
	if _, err := edr.bodyWriter.Write(b); err != nil {
		return errors.Wrap(err, "failed to write map to body writer")
	}

	buf := &bytes.Buffer{}
	if _, err := edr.bodyWriter.WriteTo(buf); err != nil {
		return errors.Wrap(err, "failed to body writer WriteTo method")
	}

	resp, err := (esapi.IndicesCreateRequest{
		Index:         edr.IndexName,
		Body:          bytes.NewReader(buf.Bytes()),
		MasterTimeout: time.Second * 5,
		Timeout:       time.Second * 5,
	}).Do(context.Background(), edr.esCli)

	return errors.Wrap(err, fmt.Sprintf("failed to call IndicesCreate, resp: %+v", resp))
}
