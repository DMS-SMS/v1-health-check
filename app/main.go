// Create package in v.1.0.0 (main)
// Most package import and dependency object generation, injection occurs in this package

package main

import (
	"github.com/docker/docker/client"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/spf13/viper"
	"log"
	"runtime"
	"time"

	"github.com/DMS-SMS/v1-health-check/app/config"
	"github.com/DMS-SMS/v1-health-check/docker"
	"github.com/DMS-SMS/v1-health-check/json"
	"github.com/DMS-SMS/v1-health-check/slack"
	_syscheckConfig "github.com/DMS-SMS/v1-health-check/syscheck/config"
	_syscheckChannelDelivery "github.com/DMS-SMS/v1-health-check/syscheck/delivery/channel"
	_syscheckRepo "github.com/DMS-SMS/v1-health-check/syscheck/repository/elasticsearch"
	_syscheckUcase "github.com/DMS-SMS/v1-health-check/syscheck/usecase"
	"github.com/DMS-SMS/v1-health-check/system"
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

	// add docker, system, slack agent
	_dkr := docker.NewAgent(dkrCli)
	_sys := system.NewAgent(dkrCli)
	_slk := slack.NewAgent(config.App.SlackAPIToken(), config.App.SlackChatChannel())

	// syscheck domain repository
	sdr := _syscheckRepo.NewESDiskCheckHistoryRepository(_syscheckConfig.App, esCli, json.MapWriter())
	scr := _syscheckRepo.NewESCPUCheckHistoryRepository(_syscheckConfig.App, esCli, json.MapWriter())
	smr := _syscheckRepo.NewESMemoryCheckHistoryRepository(_syscheckConfig.App, esCli, json.MapWriter())

	// syscheck domain usecase
	sdu := _syscheckUcase.NewDiskCheckUsecase(_syscheckConfig.App, sdr, _slk, _sys)
	scu := _syscheckUcase.NewCPUCheckUsecase(_syscheckConfig.App, scr, _slk, _sys, _dkr)
	smu := _syscheckUcase.NewMemoryCheckUsecase(_syscheckConfig.App, smr, _slk, _sys, _dkr)

	// syscheck domain delivery
	_syscheckChannelDelivery.NewDiskCheckHandler(time.Tick(_syscheckConfig.App.DiskCheckDeliveryPingCycle()), sdu)
	_syscheckChannelDelivery.NewCPUCheckHandler(time.Tick(_syscheckConfig.App.CPUCheckDeliveryPingCycle()), scu)
	_syscheckChannelDelivery.NewMemoryCheckHandler(time.Tick(_syscheckConfig.App.MemoryCheckDeliveryPingCycle()), smu)

	runtime.Goexit()
}
