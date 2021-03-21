// Create package in v.1.0.0 (main)
// Most package import and dependency object generation, injection occurs in this package

package main

import (
	"github.com/spf13/viper"
	"log"

	"github.com/DMS-SMS/v1-health-check/app/config"
)

func init() {
	// set flag to log current date, time & long file name
	log.SetFlags(log.LstdFlags | log.Llongfile)

	// set and read config file in viper package
	viper.AutomaticEnv()
	viper.SetConfigFile(config.App.ConfigFile())
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

func main() {
}
