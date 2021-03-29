// Create package in v.1.0.0 (main)
// Most package import and dependency object generation, injection occurs in this package

package main

import (
	"github.com/docker/docker/client"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
	"log"
	"runtime"
	"time"

	dockeragent "github.com/DMS-SMS/v1-health-check/agent/docker"
	slackagent "github.com/DMS-SMS/v1-health-check/agent/slack"
	sysagent "github.com/DMS-SMS/v1-health-check/agent/system"
	"github.com/DMS-SMS/v1-health-check/app/config"
	"github.com/DMS-SMS/v1-health-check/json"
	_syscheckConfig "github.com/DMS-SMS/v1-health-check/syscheck/config"
	_syscheckChannelDelivery "github.com/DMS-SMS/v1-health-check/syscheck/delivery/channel"
	_syscheckRepo "github.com/DMS-SMS/v1-health-check/syscheck/repository/elasticsearch"
	_syscheckUcase "github.com/DMS-SMS/v1-health-check/syscheck/usecase"
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
	esCli, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{config.App.ESAddress()},
	})
	if err != nil {
		log.Fatal(err)
	}

	dkrCli, err := client.NewClientWithOpts(
		client.WithVersion(config.App.DockerCliVer()),
		client.WithTimeout(config.App.DockerCliTimeout()),
	)

	_slk := slackagent.New(slack.New(config.App.SlackAPIToken()), config.App.SlackChatChannel())
	_sys := sysagent.New(dkrCli)
	_dkr := dockeragent.New(dkrCli)

	// syscheck domain repository
	sdr := _syscheckRepo.NewESDiskCheckHistoryRepository(_syscheckConfig.App, esCli, json.MapWriter())
	scr := _syscheckRepo.NewESCPUCheckHistoryRepository(_syscheckConfig.App, esCli, json.MapWriter())
	smr := _syscheckRepo.NewESMemoryCheckHistoryRepository(_syscheckConfig.App, esCli, json.MapWriter())

	// syscheck domain usecase
	sdu := _syscheckUcase.NewDiskCheckUsecase(_syscheckConfig.App, sdr, _slk, _sys)
	scu := _syscheckUcase.NewCPUCheckUsecase(_syscheckConfig.App, scr, _slk, _sys, _dkr)
	smu := _syscheckUcase.NewMemoryCheckUsecase(_syscheckConfig.App, smr, _slk, _sys, _dkr)

	// syscheck domain delivery
	_syscheckChannelDelivery.NewDiskCheckHandler(time.Tick(time.Minute * 5), sdu)
	_syscheckChannelDelivery.NewCPUCheckHandler(time.Tick(time.Minute * 5), scu)
	_syscheckChannelDelivery.NewMemoryCheckHandler(time.Tick(time.Minute * 5), smu)

	runtime.Goexit()
}
