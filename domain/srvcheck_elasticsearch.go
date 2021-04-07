// Create file in v.1.0.0
// srvcheck_elasticsearch.go is file that declare model struct & repo interface about elasticsearch check in srvcheck domain.
// also, additional method of model struct is declared in this file, too.

package domain

import (
	"context"
	"strings"
)

// ElasticsearchCheckHistory model is used for record elasticsearch check history and result
type ElasticsearchCheckHistory struct {
	// get required component by embedding serviceCheckHistoryComponent
	serviceCheckHistoryComponent

	// ActivePrimaryShards specifies active primary shards number get from elasticsearch agent
	ActivePrimaryShards int

	// ActiveShards specifies total active shards number get from elasticsearch agent
	ActiveShards int

	// UnassignedShards specifies unassigned shards number get from elasticsearch agent
	UnassignedShards int

	// ActiveShardsPercent specifies active shards percent get from elasticsearch agent
	ActiveShardsPercent float64

	// IfJaegerIndexDeleted specifies if jaeger index is deleted
	IfJaegerIndexDeleted bool

	// DeletedJaegerIndices specifies deleted jaeger indices list
	DeletedJaegerIndices []string
}

// ElasticsearchCheckHistoryRepository is interface for repository layer used in usecase layer
// Repository is implemented with elasticsearch in v.1.0.0
type ElasticsearchCheckHistoryRepository interface {
	// get required component by embedding serviceCheckHistoryRepositoryComponent
	serviceCheckHistoryRepositoryComponent

	// Store method save ElasticsearchCheckHistory model in repository
	// b in return represents bytes of response body(map[string]interface{})
	Store(*ElasticsearchCheckHistory) (b []byte, err error)
}

// ElasticsearchCheckUseCase is interface used as business process handler about elasticsearch check
type ElasticsearchCheckUseCase interface {
	// CheckElasticsearch method check elasticsearch status and store check history using repository
	CheckElasticsearch(ctx context.Context) error
}

// FillPrivateComponent overriding FillPrivateComponent method of serviceCheckHistoryComponent
func (eh *ElasticsearchCheckHistory) FillPrivateComponent() {
	eh.serviceCheckHistoryComponent.FillPrivateComponent()
	eh._type = "ElasticsearchCheck"
}

// DottedMapWithPrefix convert ElasticsearchCheckHistory to dotted map and return using MapWithPrefixKey of upper struct
// all key value of Map start with prefix received from parameter
func (eh *ElasticsearchCheckHistory) DottedMapWithPrefix(prefix string) (m map[string]interface{}) {
	m = eh.serviceCheckHistoryComponent.DottedMapWithPrefix(prefix)

	if prefix != "" {
		prefix += "."
	}

	// setting public field value in dotted map
	m[prefix + "active_primary_shards"] = eh.ActivePrimaryShards
	m[prefix + "active_shards"] = eh.ActiveShards
	m[prefix + "unassigned_shards"] = eh.UnassignedShards
	m[prefix + "active_shards_percent"] = eh.ActiveShardsPercent
	m[prefix + "if_jaeger_index_deleted"] = eh.IfJaegerIndexDeleted
	m[prefix + "deleted_jaeger_indices"] = strings.Join(eh.DeletedJaegerIndices, " | ")

	return
}

// SetClusterHealth method set field about cluster health with received cluster
func (eh *ElasticsearchCheckHistory) SetClusterHealth(cluster interface {
	ActivePrimaryShards() int     // get active primary shards number in cluster health result
	ActiveShards() int            // get active shards number in cluster health result
	UnassignedShards() int        // get unassigned shards number in cluster health result
	ActiveShardsPercent() float64 // get active shards percent in cluster health result)
}) {
	eh.ActivePrimaryShards = cluster.ActivePrimaryShards()
	eh.ActiveShards = cluster.ActiveShards()
	eh.UnassignedShards = cluster.UnassignedShards()
	eh.ActiveShardsPercent = cluster.ActiveShardsPercent()
}
