// Create package in v.1.0.0.
// config package contains App global variable with config value about srvcheck from environment variable or config file
// App return field value from method having same name with that field name

// config.go is file that define srvcheckConfig type which is type of App
// Also, App implement various config interface each of package in srvcheck domain by declaring method

package config

import (
	"github.com/spf13/viper"
	"time"
)

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

	// ---

	// fields using in elasticsearch health checking (implement elasticsearchCheckUsecaseConfig)
	// targetIndices represent indices separated with dot which are target of elasticsearch health check
	targetIndices *string

	// maximumShardsNumber represent maximum shards number of elasticsearch target cluster
	maximumShardsNumber *int

	// jaegerIndexRegexp represent jaeger index regular expression to deliver to elasticsearch agency
	jaegerIndexRegexp *string

	// jaegerIndexMinLifeCycle represent minimum life cycle of jaeger index in elasticsearch
	jaegerIndexMinLifeCycle *time.Duration
}

const (
	defaultIndexName       = "sms-service-check" // default const string for indexName
	defaultIndexShardNum   = 2                   // default const int for indexShardNum
	defaultIndexReplicaNum = 0                   // default const int for indexReplicaNum

	defaultTargetIndices           = "_all"          // default const string for TargetIndices
	defaultMaximumShardsNumber     = 900             // default const int for MaximumShardsNumber
	defaultJaegerIndexMinLifeCycle = time.Hour * 720 // default const duration for JaegerIndexMinLifeCycle
	defaultJaegerIndexRegexp       = "^jaeger-(span|service)-\\d{4}-\\d{2}-\\d{2}$" // default const string for JaegerIndexRegexp
)

// implement IndexName method of esRepositoryComponentConfig interface
func (sc *srvcheckConfig) IndexName() string {
	var key = "srvcheck.repository.elasticsearch.index.name"
	if sc.indexName == nil {
		if _, ok := viper.Get(key).(string); !ok {
			viper.Set(key, defaultIndexName)
		}
		sc.indexName = _string(viper.GetString(key))
	}
	return *sc.indexName
}

// implement IndexShardNum method of esRepositoryComponentConfig interface
func (sc *srvcheckConfig) IndexShardNum() int {
	var key = "srvcheck.repository.elasticsearch.index.shardNum"
	if sc.indexShardNum == nil {
		if _, ok := viper.Get(key).(int); !ok {
			viper.Set(key, defaultIndexShardNum)
		}
		sc.indexShardNum = _int(viper.GetInt(key))
	}
	return *sc.indexShardNum
}

// implement IndexReplicaNum method of esRepositoryComponentConfig interface
func (sc *srvcheckConfig) IndexReplicaNum() int {
	var key = "srvcheck.repository.elasticsearch.index.replicaNum"
	if sc.indexReplicaNum == nil {
		if _, ok := viper.Get(key).(int); !ok {
			viper.Set(key, defaultIndexReplicaNum)
		}
		sc.indexReplicaNum = _int(viper.GetInt(key))
	}
	return *sc.indexReplicaNum
}

// implement TargetIndices method of elasticsearchCheckUsecaseConfig interface
func (sc *srvcheckConfig) TargetIndices() string {
	var key = "srvcheck.elasticsearch.targetIndices"
	if sc.targetIndices == nil {
		if _, ok := viper.Get(key).(string); !ok {
			viper.Set(key, defaultTargetIndices)
		}
		sc.targetIndices = _string(viper.GetString(key))
	}
	return *sc.targetIndices
}

// implement MaximumShardsNumber method of elasticsearchCheckUsecaseConfig interface
func (sc *srvcheckConfig) MaximumShardsNumber() int {
	var key = "srvcheck.elasticsearch.maximumShardsNumber"
	if sc.maximumShardsNumber == nil {
		if _, ok := viper.Get(key).(int); !ok {
			viper.Set(key, defaultMaximumShardsNumber)
		}
		sc.maximumShardsNumber = _int(viper.GetInt(key))
	}
	return *sc.maximumShardsNumber
}

// implement JaegerIndexMinLifeCycle method of elasticsearchCheckUsecaseConfig interface
func (sc *srvcheckConfig) JaegerIndexMinLifeCycle() time.Duration {
	var key = "srvcheck.elasticsearch.jaegerIndexMinLifeCycle"
	if sc.jaegerIndexMinLifeCycle != nil {
		return *sc.jaegerIndexMinLifeCycle
	}

	d, err := time.ParseDuration(viper.GetString(key))
	if err != nil {
		viper.Set(key, defaultJaegerIndexMinLifeCycle.String())
		d = defaultJaegerIndexMinLifeCycle
	}

	sc.jaegerIndexMinLifeCycle = &d
	return *sc.jaegerIndexMinLifeCycle
}

// implement JaegerIndexRegexp method of elasticsearchCheckUsecaseConfig interface
func (sc *srvcheckConfig) JaegerIndexRegexp() string {
	var key = "srvcheck.elasticsearch.jaegerIndexRegexp"
	if sc.jaegerIndexRegexp == nil {
		if _, ok := viper.Get(key).(string); !ok {
			viper.Set(key, defaultJaegerIndexRegexp)
		}
		sc.jaegerIndexRegexp = _string(viper.GetString(key))
	}
	return *sc.jaegerIndexRegexp
}

// init function initialize App global variable
func init() {
	App = &srvcheckConfig{}
}

// function returns pointer variable generated from parameter
func _string(s string) *string { return &s }
func _int(i int) *int {return &i}
