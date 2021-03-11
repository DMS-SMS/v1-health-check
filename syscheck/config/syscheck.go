// Create package in v.1.0.0.
// config package contains App global variable with config value about syscheck from environment variable or config file
// App return field value from method having same name with that field name

// config.go is file that define syscheckConfig type which is type of App
// Also, App implement various config interface each of package in syscheck domain by declaring method

package config

// syscheckConfig having config value and implement various interface about Config by declaring method
type syscheckConfig struct {
	// indexName represent name of elasticsearch index including syscheck history document
	indexName *string

	// indexShardNum represent shard number of elasticsearch index storing syscheck history document
	indexShardNum *string

	// indexReplicaNum represent replica number of elasticsearch index to replace index when node become unable
	indexReplicaNum *string
}
