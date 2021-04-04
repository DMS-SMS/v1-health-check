// Create file in v.1.0.0
// srvcheck_consul_repo.go is file that define implement consul history repository using elasticsearch
// this elasticsearch repository struct embed esRepositoryRequiredComponent struct in ./srvcheck.go file

package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"

	"github.com/DMS-SMS/v1-health-check/domain"
)

// esConsulCheckHistoryRepository is to handle ConsulCheckHistoryRepository model using elasticsearch as data store
type esConsulCheckHistoryRepository struct {
	// myCfg is used for get consul history repository config about elasticsearch
	myCfg esConsulCheckHistoryRepoConfig

	// esCli is elasticsearch client connection injected from the outside package
	esCli *elasticsearch.Client

	// reqBodyWriter is implementation of reqBodyWriter interface to write []byte for request body
	reqBodyWriter reqBodyWriter
}

// esConsulCheckHistoryRepoConfig is the config for consul check history repository using elasticsearch
type esConsulCheckHistoryRepoConfig interface {
	// get common method from embedding esRepositoryComponentConfig
	esRepositoryComponentConfig
}

// NewESConsulCheckHistoryRepository return new object that implement ConsulCheckHistoryRepository interface
func NewESConsulCheckHistoryRepository(
	cfg esConsulCheckHistoryRepoConfig,
	cli *elasticsearch.Client,
	w reqBodyWriter,
) domain.ConsulCheckHistoryRepository {
	repo := &esConsulCheckHistoryRepository{
		myCfg:         cfg,
		esCli:         cli,
		reqBodyWriter: w,
	}

	if err := repo.Migrate(); err != nil {
		log.Fatal(errors.Wrap(err, "could not migrate repository").Error())
	}

	return repo
}

// Implement Migrate method of ConsulCheckHistoryRepository interface
func (ecr *esConsulCheckHistoryRepository) Migrate() error {
	resp, err := (esapi.IndicesExistsRequest{
		Index: []string{ecr.myCfg.IndexName()},
	}).Do(context.Background(), ecr.esCli)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to call IndicesExists, resp: %+v", resp))
	}

	if resp.StatusCode == http.StatusNotFound {
		if err := ecr.createIndex(); err != nil {
			return errors.Wrap(err, "failed to create index")
		}
	}

	return nil
}

// createIndex method create index with name, shard number in esRepositoryComponentConfig
func (ecr *esConsulCheckHistoryRepository) createIndex() error {
	body := map[string]interface{}{}
	body["settings.number_of_shards"] = ecr.myCfg.IndexShardNum()
	body["settings.number_of_replicas"] = ecr.myCfg.IndexReplicaNum()

	b, _ := json.Marshal(body)
	if _, err := ecr.reqBodyWriter.Write(b); err != nil {
		return errors.Wrap(err, "failed to write map to body writer")
	}

	buf := &bytes.Buffer{}
	if _, err := ecr.reqBodyWriter.WriteTo(buf); err != nil {
		return errors.Wrap(err, "failed to body writer WriteTo method")
	}

	resp, err := (esapi.IndicesCreateRequest{
		Index:         ecr.myCfg.IndexName(),
		Body:          bytes.NewReader(buf.Bytes()),
		MasterTimeout: time.Second * 5,
		Timeout:       time.Second * 5,
	}).Do(context.Background(), ecr.esCli)

	return errors.Wrap(err, fmt.Sprintf("failed to call IndicesCreate, resp: %+v", resp))
}
