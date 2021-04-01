// Create package in v.1.0.0
// elasticsearch package define struct which is implement various interface about elasticsearch agency using in usecase each of domain
// there are kind of elasticsearch agency function such as get or delete cluster, indices

// in agent.go file, define struct type of elasticsearch agent & initializer that are not method.
// Also if exist, custom type or variable used in common in each of method will declared in this file.

package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v7"
)

// slackAgent agent various elasticsearch API(get or delete cluster, indices, etc ...) as implementation
type elasticsearchAgent struct {
	// esCli is elasticsearch client connection injected from the outside package
	esCli *elasticsearch.Client
}

// NewAgent return new initialized instance of elasticsearchAgent pointer type with elasticsearch client
func NewAgent(ec *elasticsearch.Client) *elasticsearchAgent {
	return &elasticsearchAgent{
		esCli: ec,
	}
}
