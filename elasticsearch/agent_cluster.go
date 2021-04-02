// Create file in v.1.0.0
// agent_cluster.go file define method of elasticsearchAgent about cluster API
// implement agency interface about elasticsearch cluster defined in each of domain

package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DMS-SMS/v1-health-check/domain"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
	"time"
)

// GetClusterHealth return interface have various get method about cluster health inform
func (ea *elasticsearchAgent) GetClusterHealth() (interface {
	ActivePrimaryShards() int                  // get active primary shards number in cluster health result
	ActiveShards() int                         // get active shards number in cluster health result
	UnassignedShards() int                     // get unassigned shards number in cluster health result
	ActiveShardsPercent() float64              // get active shards percent in cluster health result
	WriteTo(*domain.ElasticsearchCheckHistory) // write value in result to elasticsearch check history
}, error) {
	var (
		ctx = context.Background()
	)

	resp, err := (esapi.ClusterHealthRequest{
		Index:         []string{"_all"},
		MasterTimeout: time.Second * 5,
		Timeout:       time.Second * 5,
	}).Do(ctx, ea.esCli)

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to call ClusterHealthRequest, resp: %+v", resp))
	} else if resp.IsError() {
		return nil, errors.Errorf("ClusterHealthRequest return error code, resp: %+v", resp)
	}

	m := map[string]interface{}{}
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, errors.Wrap(err, "failed to decode resp body to map")
	}

	result := getClusterHealthResult{}
	if v, ok := m["active_primary_shards"].(float64); ok {
		result.activePrimaryShards = v
	} else {
		return nil, errors.New("float64 active_primary_shards is not in resp map")
	}

	if v, ok := m["active_shards"].(float64); ok {
		result.activeShards = v
	} else {
		return nil, errors.New("float64 active_shards is not in resp map")
	}

	if v, ok := m["unassigned_shards"].(float64); ok {
		result.unassignedShards = v
	} else {
		return nil, errors.New("float64 unassigned_shards is not in resp map")
	}

	if v, ok := m["active_shards_percent_as_number"].(float64); ok {
		result.activeShardsPercent = v
	} else {
		return nil, errors.New("float64 active_shards_percent_as_number is not in resp map")
	}

	return result, nil
}

// getClusterHealthResult is implementation of GetClusterHealth return type interface
type getClusterHealthResult struct {
	activePrimaryShards, activeShards, unassignedShards, activeShardsPercent float64
}

// define return field value methods in getClusterHealthResult
func (result getClusterHealthResult) ActivePrimaryShards() int     { return int(result.activePrimaryShards) }
func (result getClusterHealthResult) ActiveShards() int            { return int(result.activeShards) }
func (result getClusterHealthResult) UnassignedShards() int        { return int(result.unassignedShards) }
func (result getClusterHealthResult) ActiveShardsPercent() float64 { return result.activeShardsPercent }

// WriteTo method write getClusterHealthResult value to history
func (result getClusterHealthResult) WriteTo(history *domain.ElasticsearchCheckHistory) {
	history.ActivePrimaryShards = result.ActivePrimaryShards()
	history.ActiveShards = result.ActiveShards()
	history.UnassignedShards = result.UnassignedShards()
	history.ActiveShardsPercent = result.ActiveShardsPercent()
}
