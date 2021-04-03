// Create package in v.1.0.0 (main)
// Most package import and dependency object generation, injection occurs in this package

package main

import (
	"github.com/docker/docker/client"
	es "github.com/elastic/go-elasticsearch/v7"
	"github.com/spf13/viper"
	"log"
	"runtime"
	"time"

	"github.com/DMS-SMS/v1-health-check/app/config"
	"github.com/DMS-SMS/v1-health-check/docker"
	"github.com/DMS-SMS/v1-health-check/elasticsearch"
	"github.com/DMS-SMS/v1-health-check/json"
	"github.com/DMS-SMS/v1-health-check/slack"
	_srvcheckConfig "github.com/DMS-SMS/v1-health-check/srvcheck/config"
	_srvcheckChannelDelivery "github.com/DMS-SMS/v1-health-check/srvcheck/delivery/channel"
	_srvcheckRepo "github.com/DMS-SMS/v1-health-check/srvcheck/repository/elasticsearch"
	_srvcheckUcase "github.com/DMS-SMS/v1-health-check/srvcheck/usecase"
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
	esCli, err := es.NewClient(es.Config{
		Addresses: []string{config.App.ESAddress()},
	})
	if err != nil {
		log.Fatal(err)
	}

	dkrCli, err := client.NewClientWithOpts(
		client.WithVersion(config.App.DockerCliVer()),
		client.WithTimeout(config.App.DockerCliTimeout()),
	)

	// add docker, system, slack, elasticsearch agent
	_dkr := docker.NewAgent(dkrCli)
	_sys := system.NewAgent(dkrCli)
	_slk := slack.NewAgent(config.App.SlackAPIToken(), config.App.SlackChatChannel())
	_esa := elasticsearch.NewAgent(esCli)

	// syscheck domain repository
	// the reason separate Repository, Usecase interface in same domain -> 서로 간의 연관성 X, 더욱 더 확실한 분리를 위해
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

	// ---

	// srvcheck domain repository
	// the reason separate Repository, Usecase interface in same domain -> 서로 간의 연관성 X, 더욱 더 확실한 분리를 위해
	ser := _srvcheckRepo.NewESElasticsearchCheckHistoryRepository(_srvcheckConfig.App, esCli, json.MapWriter())
	ssr := _srvcheckRepo.NewESSwarmpitCheckHistoryRepository(_srvcheckConfig.App, esCli, json.MapWriter())

	// srvcheck domain usecase
	seu := _srvcheckUcase.NewElasticsearchCheckUsecase(_srvcheckConfig.App, ser, _slk, _esa)
	ssu := _srvcheckUcase.NewSwarmpitCheckUsecase(_srvcheckConfig.App, ssr, _slk, _dkr)

	// srvcheck domain delivery
	_srvcheckChannelDelivery.NewElasticsearchCheckHandler(time.Tick(_srvcheckConfig.App.ESCheckDeliveryPingCycle()), seu)
	_srvcheckChannelDelivery.NewSwarmpitCheckHandler(time.Tick(_srvcheckConfig.App.SwarmpitCheckDeliveryPingCycle()), ssu)

	runtime.Goexit()
}
