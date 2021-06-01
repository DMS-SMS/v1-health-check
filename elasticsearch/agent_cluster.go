// Create file in v.1.0.0
// agent_cluster.go file define method of elasticsearchAgent about cluster API
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

// GetClusterHealth return interface have various get method about cluster health inform
func (ea *elasticsearchAgent) GetClusterHealth() (interface {
	ActivePrimaryShards() int     // get active primary shards number in cluster health result
	ActiveShards() int            // get active shards number in cluster health result
	UnassignedShards() int        // get unassigned shards number in cluster health result
	ActiveShardsPercent() float64 // get active shards percent in cluster health result
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

	var cluster cluster
	if v, ok := m["active_primary_shards"].(float64); ok {
		cluster.activePrimaryShards = v
	} else {
		return nil, errors.New("float64 active_primary_shards is not in resp map")
	}

	if v, ok := m["active_shards"].(float64); ok {
		cluster.activeShards = v
	} else {
		return nil, errors.New("float64 active_shards is not in resp map")
	}

	if v, ok := m["unassigned_shards"].(float64); ok {
		cluster.unassignedShards = v
	} else {
		return nil, errors.New("float64 unassigned_shards is not in resp map")
	}

	if v, ok := m["active_shards_percent_as_number"].(float64); ok {
		cluster.activeShardsPercent = v
	} else {
		return nil, errors.New("float64 active_shards_percent_as_number is not in resp map")
	}

	return cluster, nil
}

// cluster is struct having & handling inform about cluster, and implementation of GetClusterHealth return type interface
type cluster struct {
	activePrimaryShards, activeShards, unassignedShards, activeShardsPercent float64
}

// define return field value methods in cluster
func (c cluster) ActivePrimaryShards() int     { return int(c.activePrimaryShards) }
func (c cluster) ActiveShards() int            { return int(c.activeShards) }
func (c cluster) UnassignedShards() int        { return int(c.unassignedShards) }
func (c cluster) ActiveShardsPercent() float64 { return c.activeShardsPercent }
