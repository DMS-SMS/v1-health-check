// Create file in v.1.0.0
// agent_index.go file define method of elasticsearchAgent about index API
// implement agency interface about elasticsearch cluster defined in each of domain

package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
	"time"
)

// GetIndicesWithRegexp return indices list with regexp pattern
func (ea *elasticsearchAgent) GetIndicesWithPatterns(patterns []string) (indices []string, err error) {
	var (
		ctx = context.Background()
	)

	resp, err := (esapi.CatIndicesRequest{
		Index:         patterns,
		Format:        "JSON",
		S:             []string{"index"},
		MasterTimeout: time.Second * 5,
	}).Do(ctx, ea.esCli)

	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to call CatIndicesRequest, resp: %+v", resp))
		return
	} else if resp.IsError() {
		err = errors.Errorf("CatIndicesRequest return error code, resp: %+v", resp)
		return
	}

	var ms []map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&ms); err != nil {
		err = errors.Wrap(err, "failed to decode resp body to map slice")
		return
	}

	for _, m := range ms {
		if v, ok := m["index"].(string); ok {
			indices = append(indices, v)
		} else {
			err = errors.Wrap(err, "string index is not in resp map")
			return
		}
	}

	return
}

// DeleteIndices method delete indices in list received from parameter
func (ea *elasticsearchAgent) DeleteIndices(indices []string) (err error) {
	var (
		ctx = context.Background()
	)

	resp, err := (esapi.IndicesDeleteRequest{
		Index:         indices,
		MasterTimeout: time.Second * 5,
		Timeout:       time.Second * 5,
	}).Do(ctx, ea.esCli)

	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to call CatIndicesRequest, resp: %+v", resp))
	} else if resp.IsError() {
		err = errors.Errorf("CatIndicesRequest return error code, resp: %+v", resp)
	}
	return
}

// indices is struct having & handling indices inform, and implementation of GetIndicesWithPatterns return type interface
type indices []struct {
	name    string    // specifies index name
	created time.Time // specifies index created time
}
