// Create package in v.1.0.0.
// config package contains App global variable with config value about srvcheck from environment variable or config file
// App return field value from method having same name with that field name

// config.go is file that define srvcheckConfig type which is type of App
// Also, App implement various config interface each of package in srvcheck domain by declaring method

package config

// App is the application config about srvcheck domain
var App *srvcheckConfig

// srvcheckConfig having config value and implement various interface about Config by declaring method
type srvcheckConfig struct {
	// fields about index information in elasticsearch (implement esRepositoryComponentConfig)
	// indexName represent name of elasticsearch index including srvcheck history document
	indexName *string

	// indexShardNum represent shard number of elasticsearch index storing srvcheck history document
	indexShardNum *int

	// indexReplicaNum represent replica number of elasticsearch index to replace index when node become unable
	indexReplicaNum *int
}

// init function initialize App global variable
func init() {
	App = &srvcheckConfig{}
}
