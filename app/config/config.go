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

	// awsAccountID represent aws account ID
	awsAccountID *string

	// awsAccountKey represent aws account key
	awsAccountKey *string

	// awsRegion represent aws region
	awsRegion *string

	// awsS3Bucket represent aws s3 bucket
	awsS3Bucket *string

	// version represent version of sms health check(this application)
	version *string
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

// Version return version from environment variable
func (ac *appConfig) Version() string {
	if ac.version != nil {
		return *ac.version
	}

	if viper.IsSet("VERSION") {
		ac.version = _string(viper.GetString("VERSION"))
	} else {
		log.Fatal("please set VERSION in environment variable")
	}
	return *ac.version
}

// return docker client version as literal
func (ac *appConfig) DockerCliVer() string {
	return "1.40"
}

// return docker client connection time out as literal
func (ac *appConfig) DockerCliTimeout() time.Duration {
	return time.Second * 5
}

// AWSS3Bucket return s3 bucket name from environment variable
func (ac *appConfig) AWSS3Bucket() string {
	if ac.awsS3Bucket != nil {
		return *ac.awsS3Bucket
	}

	if viper.IsSet("SMS_AWS_BUCKET") {
		ac.awsS3Bucket = _string(viper.GetString("SMS_AWS_BUCKET"))
	} else {
		log.Fatal("please set SMS_AWS_BUCKET in environment variable")
	}
	return *ac.awsS3Bucket
}

func (ac *appConfig) ProfilerS3Bucket() string {
	return ac.AWSS3Bucket()
}

// AWSAccountID return aws account id from environment variable
func (ac *appConfig) AWSAccountID() string {
	if ac.awsAccountID != nil {
		return *ac.awsAccountID
	}

	if viper.IsSet("SMS_AWS_ID") {
		ac.awsAccountID = _string(viper.GetString("SMS_AWS_ID"))
	} else {
		log.Fatal("please set SMS_AWS_ID in environment variable")
	}
	return *ac.awsAccountID
}

// AWSAccountKey return aws account key from environment variable
func (ac *appConfig) AWSAccountKey() string {
	if ac.awsAccountKey != nil {
		return *ac.awsAccountKey
	}

	if viper.IsSet("SMS_AWS_KEY") {
		ac.awsAccountKey = _string(viper.GetString("SMS_AWS_KEY"))
	} else {
		log.Fatal("please set SMS_AWS_KEY in environment variable")
	}
	return *ac.awsAccountKey
}

// AWSRegion return aws region from environment variable
func (ac *appConfig) AWSRegion() string {
	if ac.awsRegion != nil {
		return *ac.awsRegion
	}

	if viper.IsSet("SMS_AWS_REGION") {
		ac.awsRegion = _string(viper.GetString("SMS_AWS_REGION"))
	} else {
		log.Fatal("please set SMS_AWS_REGION in environment variable")
	}
	return *ac.awsRegion
}

func init() {
	App = &appConfig{}
}

// _string and _int function returns pointer variable generated from parameter
func _string(s string) *string { return &s }
