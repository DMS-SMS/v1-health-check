// Create package in v.1.0.0
// config package contains App global variable with config value using in app(main) package from environment variable or config file
// App return field value from method having same name with that field name

// config.go is file that define appConfig type which is type of App
// appConfig dose not implement any interface, so use by explicitly importing in app(main) package
// Also, App implement various config interface each of package in syscheck domain by declaring method

package config

import (
	"log"
	"os"
	"time"
)

// App is the application config using in main package
var App *appConfig

// appConfig having config value and return that value with method. Not implement interface
type appConfig struct {
	// esAddress represent host address of elasticsearch server
	esAddress *string

	// configFile represent full name of config file
	configFile *string
}

// return elasticsearch address get from environment variable
func (ac *appConfig) ESAddress() string {
	if ac.esAddress != nil {
		return *ac.esAddress
	}

	if v := os.Getenv("ES_ADDRESS"); v == "" {
		log.Fatal("please set ES_ADDRESS in environment variable")
	} else {
		ac.esAddress = &v
	}
	return *ac.esAddress
}

// return elasticsearch address get from environment variable
func (ac *appConfig) ConfigFile() string {
	if ac.configFile != nil {
		return *ac.configFile
	}

	if v := os.Getenv("CONFIG_FILE"); v == "" {
		log.Fatal("please set CONFIG_FILE in environment variable")
	} else {
		ac.configFile = &v
	}
	return *ac.configFile
}

func (ac *appConfig) DockerCliVer() string {
	return "1.40"
}

func (ac *appConfig) DockerCliTimeout() time.Duration {
	return time.Second * 5
}

func init() {
	App = &appConfig{}
}
