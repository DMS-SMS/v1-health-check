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
func (ea *elasticsearchAgent) GetIndicesWithPatterns(patterns []string) (interface {
	SetMinLifeCycle(cycle time.Duration) // set min life cycle of index of indices
	IndexNames() []string                // get index name list of indices
}, error) {
	var (
		ctx = context.Background()
	)

	resp, err := (esapi.CatIndicesRequest{
		Index:         patterns,
		Format:        "JSON",
		MasterTimeout: time.Second * 5,

		S: []string{"index"},
		H: []string{"index", "creation.date.string"},
	}).Do(ctx, ea.esCli)

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to call CatIndicesRequest, resp: %+v", resp))
	} else if resp.IsError() {
		return nil, errors.Errorf("CatIndicesRequest return error code, resp: %+v", resp)
	}

	var ms []map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&ms); err != nil {
		return nil, errors.Wrap(err, "failed to decode resp body to map slice")
	}

	indices := make(indices, len(ms))
	var idx = 0
	for _, m := range ms {
		if v, ok := m["index"].(string); ok {
			indices[idx].name = v
		} else {
			return nil, errors.Wrap(err, "string index key is not in resp map")
		}

		if v, ok := m["creation.date.string"].(string); ok {
			t, _ := time.Parse(time.RFC3339, v)
			indices[idx].created = t
		} else {
			return nil, errors.Wrap(err, "string creation.date.string key is not in resp map")
		}
		idx++
	}

	return &indices, nil
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

// SetMinLifeCycle set minimum life cycle of index & reset value of receiver variable
func (idxes *indices) SetMinLifeCycle(cycle time.Duration) {
	life := time.Now().Add(-cycle)

	var filtered indices
	for _, index := range *idxes {
		if index.created.Before(life) {
			filtered = append(filtered, index)
		}
	}
	*idxes = filtered
}

// IndexNames return index name list of indices
func (idxes indices) IndexNames() (names []string) {
	names = make([]string, len(idxes))
	for i, index := range idxes {
		names[i] = index.name
	}
	return
}
