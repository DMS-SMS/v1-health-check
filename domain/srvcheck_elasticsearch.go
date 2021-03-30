// Create file in v.1.0.0
// srvcheck_elasticsearch.go is file that declare model struct & repo interface about elasticsearch check in srvcheck domain.
// also, additional method of model struct is declared in this file, too.

package domain

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
