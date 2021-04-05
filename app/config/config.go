// Create package in v.1.0.0
// config package contains App global variable with config value using in app(main) package from environment variable or config file
// App return field value from method having same name with that field name

// config.go is file that define appConfig type which is type of App
// appConfig dose not implement any interface, so use by explicitly importing in app(main) package
// Also, App implement various config interface each of package in syscheck domain by declaring method

package config

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

// App is the application config using in main package
var App *appConfig

// appConfig having config value and return that value with method. Not implement interface
type appConfig struct {
	// esAddress represent host address of elasticsearch server
	esAddress *string

	// consulAddress represent host address of consul server
	consulAddress *string

	// configFile represent full name of config file
	configFile *string

	// slackAPIToken represent token to using in slack API
	slackAPIToken *string

	// slackChatCnl represent slack channel ID to send chat
	slackChatCnl *string
}

// return elasticsearch address get from environment variable
func (ac *appConfig) ESAddress() string {
	if ac.esAddress != nil {
		return *ac.esAddress
	}

	if viper.IsSet("ES_ADDRESS") {
		ac.esAddress = _string(viper.GetString("ES_ADDRESS"))
	} else {
		log.Fatal("please set ES_ADDRESS in environment variable")
	}
	return *ac.esAddress
}

// return consul address get from environment variable
func (ac *appConfig) ConsulAddress() string {
	if ac.consulAddress != nil {
		return *ac.consulAddress
	}

	if viper.IsSet("CONSUL_ADDRESS") {
		ac.consulAddress = _string(viper.GetString("CONSUL_ADDRESS"))
	} else {
		log.Fatal("please set CONSUL_ADDRESS in environment variable")
	}
	return *ac.consulAddress
}

// return elasticsearch address get from environment variable
func (ac *appConfig) ConfigFile() string {
	if ac.configFile != nil {
		return *ac.configFile
	}

	if viper.IsSet("CONFIG_FILE") {
		ac.configFile = _string(viper.GetString("CONFIG_FILE"))
	} else {
		log.Fatal("please set CONFIG_FILE in environment variable")
	}
	return *ac.configFile
}

// return slack api token get from environment variable
func (ac *appConfig) SlackAPIToken() string {
	if ac.slackAPIToken != nil {
		return *ac.slackAPIToken
	}

	if viper.IsSet("SLACK_API_TOKEN") {
		ac.slackAPIToken = _string(viper.GetString("SLACK_API_TOKEN"))
	} else {
		log.Fatal("please set SLACK_API_TOKEN in environment variable")
	}
	return *ac.slackAPIToken
}

// return slack chat channel ID get from environment variable
func (ac *appConfig) SlackChatChannel() string {
	if ac.slackChatCnl != nil {
		return *ac.slackChatCnl
	}

	if viper.IsSet("SLACK_CHAT_CHANNEL") {
		ac.slackChatCnl = _string(viper.GetString("SLACK_CHAT_CHANNEL"))
	} else {
		log.Fatal("please set SLACK_CHAT_CHANNEL in environment variable")
	}
	return *ac.slackChatCnl
}

// return docker client version as literal
func (ac *appConfig) DockerCliVer() string {
	return "1.40"
}

// return docker client connection time out as literal
func (ac *appConfig) DockerCliTimeout() time.Duration {
	return time.Second * 5
}

func init() {
	App = &appConfig{}
}

// _string and _int function returns pointer variable generated from parameter
func _string(s string) *string { return &s }
