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
)

// App is the application config using in main package
var App *appConfig

// appConfig having config value and return that value with method. Not implement interface
type appConfig struct {
	// esAddress represent host address of elasticsearch server
	esAddress *string
}

// return elasticsearch address get from environment variable
func (ac *appConfig) ESAddress() string {
	if ac.esAddress != nil {
		return *ac.esAddress
	}

	esAddr := os.Getenv("ES_ADDRESS")
	if esAddr == "" {
		log.Fatal("please set ES_ADDRESS in environment variable")
	}

	ac.esAddress = &esAddr
	return *ac.esAddress
}

func init() {
	App = &appConfig{}
}
