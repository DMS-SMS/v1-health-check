// Create file in v.1.0.0
// syscheck_disk_repo.go is file that define repository implement about disk using elasticsearch
// disk repository struct embed esRepositoryRequiredComponent struct in ./syscheck.go file

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

// esDiskCheckHistoryRepository is to handle DiskCheckHistory model using elasticsearch as data store
type esDiskCheckHistoryRepository struct {
	// myCfg is used for get disk check history repository config about elasticsearch
	myCfg esDiskCheckHistoryRepoConfig

	// esCli is elasticsearch client connection injected from the outside package
	esCli *elasticsearch.Client

	// bodyWriter is implementation of reqBodyWriter interface to write []byte for request body
	bodyWriter reqBodyWriter
}

// esDiskCheckHistoryRepoConfig is the config for disk check history repository using elasticsearch
type esDiskCheckHistoryRepoConfig interface {
	// get common method from embedding esRepositoryComponentConfig
	esRepositoryComponentConfig
}

// NewESDiskCheckHistoryRepository return new object that implement DiskCheckHistory.Repository interface
func NewESDiskCheckHistoryRepository(cfg esDiskCheckHistoryRepoConfig, cli *elasticsearch.Client, w reqBodyWriter) domain.DiskCheckHistoryRepository {
	repo := &esDiskCheckHistoryRepository{
		myCfg:      cfg,
		esCli:      cli,
		bodyWriter: w,
	}

	if err := repo.Migrate(); err != nil {
		log.Fatal(errors.Wrap(err, "could not migrate repository").Error())
	}

	return repo
}

// Implement Migrate method of DiskCheckHistoryRepository interface
func (edr *esDiskCheckHistoryRepository) Migrate() error {
	resp, err := (esapi.IndicesExistsRequest{
		Index: []string{edr.myCfg.IndexName()},
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
	body["settings.number_of_shards"] = edr.myCfg.IndexShardNum()
	body["settings.number_of_replicas"] = edr.myCfg.IndexReplicaNum()

	b, _ := json.Marshal(body)
	if _, err := edr.bodyWriter.Write(b); err != nil {
		return errors.Wrap(err, "failed to write map to body writer")
	}

	buf := &bytes.Buffer{}
	if _, err := edr.bodyWriter.WriteTo(buf); err != nil {
		return errors.Wrap(err, "failed to body writer WriteTo method")
	}

	resp, err := (esapi.IndicesCreateRequest{
		Index:         edr.myCfg.IndexName(),
		Body:          bytes.NewReader(buf.Bytes()),
		MasterTimeout: time.Second * 5,
		Timeout:       time.Second * 5,
	}).Do(context.Background(), edr.esCli)

	return errors.Wrap(err, fmt.Sprintf("failed to call IndicesCreate, resp: %+v", resp))
}

// Implement Store method of DiskCheckHistoryRepository interface
func (edr *esDiskCheckHistoryRepository) Store(history *domain.DiskCheckHistory) (b []byte, err error) {
	body, _ := json.Marshal(history.DottedMapWithPrefix(""))
	if _, err = edr.bodyWriter.Write(body); err != nil {
		err = errors.Wrap(err, "failed to write map to body writer")
		return
	}

	buf := &bytes.Buffer{}
	if _, err = edr.bodyWriter.WriteTo(buf); err != nil {
		err = errors.Wrap(err, "failed to body writer WriteTo method")
		return
	}

	resp, err := (esapi.IndexRequest{
		Index:        edr.myCfg.IndexName(),
		Body:         bytes.NewReader(buf.Bytes()),
		Timeout:      time.Second * 5,
	}).Do(context.Background(), edr.esCli)

	if err != nil || resp.IsError() {
		// 단순 err인지 status code상 오류인지 확인
		err = errors.Wrap(err, fmt.Sprintf("failed to call IndexRequest, resp: %+v", resp))
		return
	}

	result := map[string]interface{}{}
	_ = json.NewDecoder(resp.Body).Decode(&result)
	b, _ = json.Marshal(result)
	return
}
