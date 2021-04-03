// Create file in v.1.0.0
// syscheck_memory_repo.go is file that define repository implement about memory using elasticsearch
// memory repository struct embed esRepositoryRequiredComponent struct in ./syscheck.go file

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

// Implement Migrate method of MemoryCheckHistoryRepository interface
func (emr *esMemoryCheckHistoryRepository) Migrate() error {
	resp, err := (esapi.IndicesExistsRequest{
		Index: []string{emr.myCfg.IndexName()},
	}).Do(context.Background(), emr.esCli)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to call IndicesExists, resp: %+v", resp))
	}

	if resp.StatusCode == http.StatusNotFound {
		if err := emr.createIndex(); err != nil {
			return errors.Wrap(err, "failed to create index")
		}
	}

	return nil
}

// createIndex method create index with name, share number in esMemoryCheckHistoryRepository
func (emr *esMemoryCheckHistoryRepository) createIndex() error {
	body := map[string]interface{}{}
	body["settings.number_of_shards"] = emr.myCfg.IndexShardNum()
	body["settings.number_of_replicas"] = emr.myCfg.IndexReplicaNum()

	b, _ := json.Marshal(body)
	if _, err := emr.bodyWriter.Write(b); err != nil {
		return errors.Wrap(err, "failed to write map to body writer")
	}

	buf := &bytes.Buffer{}
	if _, err := emr.bodyWriter.WriteTo(buf); err != nil {
		return errors.Wrap(err, "failed to body writer WriteTo method")
	}

	resp, err := (esapi.IndicesCreateRequest{
		Index:         emr.myCfg.IndexName(),
		Body:          bytes.NewReader(buf.Bytes()),
		MasterTimeout: time.Second * 5,
		Timeout:       time.Second * 5,
	}).Do(context.Background(), emr.esCli)

	return errors.Wrap(err, fmt.Sprintf("failed to call IndicesCreate, resp: %+v", resp))
}

// Implement Store method of MemoryCheckHistoryRepository interface
func (emr *esMemoryCheckHistoryRepository) Store(history *domain.MemoryCheckHistory) (b []byte, err error) {
	body, _ := json.Marshal(history.DottedMapWithPrefix(""))
	if _, err = emr.bodyWriter.Write(body); err != nil {
		err = errors.Wrap(err, "failed to write map to body writer")
		return
	}

	buf := &bytes.Buffer{}
	if _, err = emr.bodyWriter.WriteTo(buf); err != nil {
		err = errors.Wrap(err, "failed to body writer WriteTo method")
		return
	}

	resp, err := (esapi.IndexRequest{
		Index:        emr.myCfg.IndexName(),
		Body:         bytes.NewReader(buf.Bytes()),
		Timeout:      time.Second * 5,
	}).Do(context.Background(), emr.esCli)

	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to call IndexRequest, resp: %+v", resp))
		return
	} else if resp.IsError() {
		err = errors.Errorf("IndexRequest return error code, resp: %+v", resp)
		return
	}

	result := map[string]interface{}{}
	_ = json.NewDecoder(resp.Body).Decode(&result)
	b, _ = json.Marshal(result)
	return
}
