// Create package in v.1.0.0.
// config package contains App global variable with config value about syscheck from environment variable or config file
// App return field value from method having same name with that field name

// config.go is file that define syscheckConfig type which is type of App
// Also, App implement various config interface each of package in syscheck domain by declaring method

package config

import (
	"github.com/inhies/go-bytesize"
	"github.com/spf13/viper"
	"time"
)

// App is the application config about syscheck domain
var App *syscheckConfig

// syscheckConfig having config value and implement various interface about Config by declaring method
type syscheckConfig struct {
	// fields about index information in elasticsearch (implement esRepositoryComponentConfig)
	// indexName represent name of elasticsearch index including syscheck history document
	indexName *string

	// indexShardNum represent shard number of elasticsearch index storing syscheck history document
	indexShardNum *int

	// indexReplicaNum represent replica number of elasticsearch index to replace index when node become unable
	indexReplicaNum *int

	// ---

	// fields using in disk health checking (implement diskCheckUsecaseConfig)
	// diskMinCapacity represent minimum disk capacity and is standard to decide to if disk is healthy.
	diskMinCapacity *bytesize.ByteSize

	// ---

	// fields using in cpu health check (implement cpuCheckUsecaseConfig)
	// cpuWarningUsage represent warning cpu usage.
	cpuWarningUsage *float64

	// cpuMaximumUsage represent cpu maximum usage and is standard to decide to if cpu is healthy.
	cpuMaximumUsage *float64

	// cpuMinimumUsageToRemove represent minimum cpu usage to decide whether remove container or not
	cpuMinimumUsageToRemove *float64

	// ---

	// fields using in memory health check (implement memoryCheckUsecaseConfig)
	// memoryWarningUsage represent warning memory usage.
	memoryWarningUsage *bytesize.ByteSize

	// memoryMaximumUsage represent memory maximum usage and is standard to decide to if memory is healthy.
	memoryMaximumUsage *bytesize.ByteSize

	// memoryMinimumUsageToRemove represent minimum memory usage to decide whether remove container or not
	memoryMinimumUsageToRemove *bytesize.ByteSize

	// --

	// fields using in main function to inject delivery layer (not implement any interface)
	// diskCheckChannelPingCycle represent disk check delivery ping cycle
	diskCheckDeliveryPingCycle *time.Duration

	// cpuCheckDeliveryPingCycle represent cpu check delivery ping cycle
	cpuCheckDeliveryPingCycle *time.Duration

	// memoryCheckDeliveryPingCycle represent memory check delivery ping cycle
	memoryCheckDeliveryPingCycle *time.Duration
}

// default const value about syscheckConfig field
const (
	defaultIndexName       = "sms-system-check" // default const string for indexName
	defaultIndexShardNum   = 2                  // default const int for indexShardNum
	defaultIndexReplicaNum = 0                  // default const int for indexReplicaNum

	defaultDiskMinCapacity = bytesize.GB * 2 // default const byte size for diskMinCapacity

	defaultCPUWarningUsage         = float64(1.0) // default const float64 for cpuWarningUsage
	defaultCPUMaximumUsage         = float64(1.5) // default const float64 for cpuMaximumUsage
	defaultCPUMinimumUsageToRemove = float64(0.5) // default const float64 for cpuMinimumUsageToRemove

	defaultMemoryWarningUsage         = bytesize.GB * 6 // default const float64 for memoryWarningUsage
	defaultMemoryMaximumUsage         = bytesize.GB * 7 // default const float64 for memoryMaximumUsage
	defaultMemoryMinimumUsageToRemove = bytesize.GB * 1 // default const float64 for memoryMinimumUsageToRemove
)

// implement IndexName method of esRepositoryComponentConfig interface
func (sc *syscheckConfig) IndexName() string {
	var key = "domain.syscheck.repository.elasticsearch.index.name"
	if sc.indexName == nil {
		if _, ok := viper.Get(key).(string); !ok {
			viper.Set(key, defaultIndexName)
		}
		sc.indexName = _string(viper.GetString(key))
	}
	return *sc.indexName
}

// implement IndexShardNum method of esRepositoryComponentConfig interface
func (sc *syscheckConfig) IndexShardNum() int {
	var key = "domain.syscheck.repository.elasticsearch.index.shardNum"
	if sc.indexShardNum == nil {
		if _, ok := viper.Get(key).(int); !ok {
			viper.Set(key, defaultIndexShardNum)
		}
		sc.indexShardNum = _int(viper.GetInt(key))
	}
	return *sc.indexShardNum
}

// implement IndexReplicaNum method of esRepositoryComponentConfig interface
func (sc *syscheckConfig) IndexReplicaNum() int {
	var key = "domain.syscheck.repository.elasticsearch.index.replicaNum"
	if sc.indexReplicaNum == nil {
		if _, ok := viper.Get(key).(int); !ok {
			viper.Set(key, defaultIndexReplicaNum)
		}
		sc.indexReplicaNum = _int(viper.GetInt(key))
	}
	return *sc.indexReplicaNum
}

// implement DiskMinCapacity method of diskCheckUsecaseConfig interface
func (sc *syscheckConfig) DiskMinCapacity() bytesize.ByteSize {
	var key = "domain.syscheck.diskcheck.minCapacity"
	if sc.diskMinCapacity != nil {
		return *sc.diskMinCapacity
	}

	size, err := bytesize.Parse(viper.GetString(key))
	if err != nil {
		viper.Set(key, defaultDiskMinCapacity.String())
		size = defaultDiskMinCapacity
	}

	sc.diskMinCapacity = &size
	return *sc.diskMinCapacity
}

// implement CPUWarningUsage method of cpuCheckUsecaseConfig interface
func (sc *syscheckConfig) CPUWarningUsage() float64 {
	var key = "domain.syscheck.cpucheck.cpuWarningUsage"
	if sc.cpuWarningUsage == nil {
		if _, ok := viper.Get(key).(float64); !ok {
			viper.Set(key, defaultCPUWarningUsage)
		}
		sc.cpuWarningUsage = _float64(viper.GetFloat64(key))
	}
	return *sc.cpuWarningUsage
}

// implement CPUMaximumUsage method of cpuCheckUsecaseConfig interface
func (sc *syscheckConfig) CPUMaximumUsage() float64 {
	var key = "domain.syscheck.cpucheck.cpuMaximumUsage"
	if sc.cpuMaximumUsage == nil {
		if _, ok := viper.Get(key).(float64); !ok {
			viper.Set(key, defaultCPUMaximumUsage)
		}
		sc.cpuMaximumUsage = _float64(viper.GetFloat64(key))
	}
	return *sc.cpuMaximumUsage
}

// CPUMinimumUsageToRemove method returns float64 represent cpu minimum usage to remove
func (sc *syscheckConfig) CPUMinimumUsageToRemove() float64 {
	var key = "domain.syscheck.cpucheck.cpuMinimumUsageToRemove"
	if sc.cpuMinimumUsageToRemove == nil {
		if _, ok := viper.Get(key).(float64); !ok {
			viper.Set(key, defaultCPUMinimumUsageToRemove)
		}
		sc.cpuMinimumUsageToRemove = _float64(viper.GetFloat64(key))
	}
	return *sc.cpuMinimumUsageToRemove
}

// implement MemoryWarningUsage method of memoryCheckUsecaseConfig interface
func (sc *syscheckConfig) MemoryWarningUsage() bytesize.ByteSize {
	var key = "domain.syscheck.memorycheck.memoryWarningUsage"
	if sc.memoryWarningUsage != nil {
		return *sc.memoryWarningUsage
	}

	size, err := bytesize.Parse(viper.GetString(key))
	if err != nil {
		viper.Set(key, defaultMemoryWarningUsage.String())
		size = defaultMemoryWarningUsage
	}

	sc.memoryWarningUsage = &size
	return *sc.memoryWarningUsage
}

// implement MemoryMaximumUsage method of memoryCheckUsecaseConfig interface
func (sc *syscheckConfig) MemoryMaximumUsage() bytesize.ByteSize {
	var key = "domain.syscheck.memorycheck.memoryMaximumUsage"
	if sc.memoryMaximumUsage != nil {
		return *sc.memoryMaximumUsage
	}

	size, err := bytesize.Parse(viper.GetString(key))
	if err != nil {
		viper.Set(key, defaultMemoryMaximumUsage.String())
		size = defaultMemoryMaximumUsage
	}

	sc.memoryMaximumUsage = &size
	return *sc.memoryMaximumUsage
}

// implement MemoryMinimumUsageToRemove method of memoryCheckUsecaseConfig interface
func (sc *syscheckConfig) MemoryMinimumUsageToRemove() bytesize.ByteSize {
	var key = "domain.syscheck.memorycheck.memoryMinimumUsageToRemove"
	if sc.memoryMinimumUsageToRemove != nil {
		return *sc.memoryMinimumUsageToRemove
	}

	size, err := bytesize.Parse(viper.GetString(key))
	if err != nil {
		viper.Set(key, defaultMemoryMinimumUsageToRemove.String())
		size = defaultMemoryMinimumUsageToRemove
	}

	sc.memoryMinimumUsageToRemove = &size
	return *sc.memoryMinimumUsageToRemove
}

// init function initialize App global variable
func init() {
	App = &syscheckConfig{}
}

// function returns pointer variable generated from parameter
func _string(s string) *string { return &s }
func _int(i int) *int {return &i}
func _float64(f float64) *float64 {return &f}
