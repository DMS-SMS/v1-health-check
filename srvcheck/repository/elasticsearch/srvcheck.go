// Create package in v.1.0.0
// elasticsearch package is for implementations of srvcheck domain repository using elasticsearch
// In practice, repository struct declaration and implementation occur in this package

// srvcheck.go is file that define structure to embed from another structures.
// It also defines interface or function used jointly in the package as private.

package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"time"
)

// esRepositoryComponentConfig is interface contains method to return config value that elasticsearch repository should have
// It can be externally set as Config object that implements that interface.
type esRepositoryComponentConfig interface {
	// IndexName method returns the index name of elasticsearch about srvcheck
	IndexName() string

	// IndexShardNum method returns the number of index shard in elasticsearch about srvcheck
	IndexShardNum() int

	// IndexReplicaNum method returns the number of index replica in elasticsearch about srvcheck
	IndexReplicaNum() int
}

// reqBodyWriter is private interface to use as writing []byte for request body
type reqBodyWriter interface {
	io.Writer
	io.WriterTo
}

// esRepositoryMigrator is struct that Migrate es repository using parameter variable
type esRepositoryMigrator struct{}

// Migrate method, if index doesn't exist, create index with name and shard number in esRepositoryComponentConfig
func (erm esRepositoryMigrator) Migrate(cfg esRepositoryComponentConfig, cli *elasticsearch.Client, w reqBodyWriter) error {
	resp, err := (esapi.IndicesExistsRequest{
		Index: []string{cfg.IndexName()},
	}).Do(context.Background(), cli)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to call IndicesExists, resp: %+v", resp))
	}

	if resp.StatusCode == http.StatusNotFound {
		body := map[string]interface{}{}
		body["settings.number_of_shards"] = cfg.IndexShardNum()
		body["settings.number_of_replicas"] = cfg.IndexReplicaNum()

		b, _ := json.Marshal(body)
		if _, err := w.Write(b); err != nil {
			return errors.Wrap(err, "failed to write map to body writer")
		}

		buf := &bytes.Buffer{}
		if _, err := w.WriteTo(buf); err != nil {
			return errors.Wrap(err, "failed to body writer WriteTo method")
		}

		if resp, err := (esapi.IndicesCreateRequest{
			Index:         cfg.IndexName(),
			Body:          bytes.NewReader(buf.Bytes()),
			MasterTimeout: time.Second * 5,
			Timeout:       time.Second * 5,
		}).Do(context.Background(), cli); err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to call IndicesCreate, resp: %+v", resp))
		}
	}

	return nil
}
