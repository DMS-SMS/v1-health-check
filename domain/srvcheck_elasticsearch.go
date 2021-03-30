// Create file in v.1.0.0
// srvcheck_elasticsearch.go is file that declare model struct & repo interface about elasticsearch check in srvcheck domain.
// also, additional method of model struct is declared in this file, too.

package domain

import "context"

// ElasticsearchCheckHistory model is used for record elasticsearch check history and result
type ElasticsearchCheckHistory struct {
	// get required component by embedding serviceCheckHistoryComponent
	serviceCheckHistoryComponent

	// activePrimaryShards specifies active primary shards number get from elasticsearch agent
	activePrimaryShards int

	// activeShards specifies total active shards number get from elasticsearch agent
	activeShards int

	// unassignedShards specifies unassigned shards number get from elasticsearch agent
	unassignedShards int

	// activeShardsPercent specifies active shards percent get from elasticsearch agent
	activeShardsPercent int
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
	m[prefix + "active_primary_shards"] = eh.activePrimaryShards
	m[prefix + "active_shard"] = eh.activeShards
	m[prefix + "unassigned_shards"] = eh.unassignedShards
	m[prefix + "active_shards_percent"] = eh.activeShardsPercent

	return
}
