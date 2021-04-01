// Create file in v.1.0.0
// agent_cluster.go file define method of elasticsearchAgent about cluster API
// implement agency interface about elasticsearch cluster defined in each of domain

package elasticsearch


// getClusterHealthResult is implementation of GetClusterHealth return type interface
type getClusterHealthResult struct {
	activePrimaryShards, activeShards, unassignedShards, activeShardsPercent float64
}

// define getClusterHealthResult methods that return field of this struct
func (result getClusterHealthResult) ActivePrimaryShards() int     { return int(result.activePrimaryShards) }
func (result getClusterHealthResult) ActiveShards() int            { return int(result.activeShards) }
func (result getClusterHealthResult) UnassignedShards() int        { return int(result.unassignedShards) }
func (result getClusterHealthResult) ActiveShardsPercent() float64 { return result.activeShardsPercent }
