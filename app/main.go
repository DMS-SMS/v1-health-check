// Create package in v.1.0.0 (main)
// Most package import and dependency object generation, injection occurs in this package

package main

import (
	// import Go SDK package
	"log"
	"runtime"
	"time"

	// import external package
	"github.com/docker/docker/client"
	es "github.com/elastic/go-elasticsearch/v7"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	// import app config & various agent package
	"github.com/DMS-SMS/v1-health-check/app/config"
	"github.com/DMS-SMS/v1-health-check/consul"
	"github.com/DMS-SMS/v1-health-check/docker"
	"github.com/DMS-SMS/v1-health-check/elasticsearch"
	"github.com/DMS-SMS/v1-health-check/json"
	"github.com/DMS-SMS/v1-health-check/slack"
	"github.com/DMS-SMS/v1-health-check/system"

	// import system check domain package
	_syscheckConfig       "github.com/DMS-SMS/v1-health-check/syscheck/config"
	_syscheckChanDelivery "github.com/DMS-SMS/v1-health-check/syscheck/delivery/channel"
	_syscheckRepo         "github.com/DMS-SMS/v1-health-check/syscheck/repository/elasticsearch"
	_syscheckUcase        "github.com/DMS-SMS/v1-health-check/syscheck/usecase"

	// import service check domain package
	_srvcheckConfig       "github.com/DMS-SMS/v1-health-check/srvcheck/config"
	_srvcheckChanDelivery "github.com/DMS-SMS/v1-health-check/srvcheck/delivery/channel"
	_srvcheckRepo         "github.com/DMS-SMS/v1-health-check/srvcheck/repository/elasticsearch"
	_srvcheckUcase        "github.com/DMS-SMS/v1-health-check/srvcheck/usecase"
)

func init() {
	// set flag to log current date, time & long file name
	log.SetFlags(log.LstdFlags | log.Lshortfile)

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
		log.Fatal(errors.Wrap(err, "failed to create elasticsearch client"))
	}

	dkrCli, err := client.NewClientWithOpts(
		client.WithVersion(config.App.DockerCliVer()),
		client.WithTimeout(config.App.DockerCliTimeout()),
	)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create docker client"))
	}

	cslCfg := api.DefaultConfig()
	cslCfg.Address = config.App.ConsulAddress()
	cslCli, err := api.NewClient(cslCfg)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create consul client"))
	}

	// add docker, system, slack, elasticsearch agent
	_dkr := docker.NewAgent(dkrCli)
	_sys := system.NewAgent(dkrCli)
	_slk := slack.NewAgent(config.App.SlackAPIToken(), config.App.SlackChatChannel())
	_es := elasticsearch.NewAgent(esCli)
	_csl := consul.NewAgent(cslCli)

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
	_syscheckChanDelivery.NewDiskCheckHandler(time.Tick(_syscheckConfig.App.DiskCheckDeliveryPingCycle()), sdu)
	_syscheckChanDelivery.NewCPUCheckHandler(time.Tick(_syscheckConfig.App.CPUCheckDeliveryPingCycle()), scu)
	_syscheckChanDelivery.NewMemoryCheckHandler(time.Tick(_syscheckConfig.App.MemoryCheckDeliveryPingCycle()), smu)

	// ---

	// srvcheck domain repository
	// the reason separate Repository, Usecase interface in same domain -> 서로 간의 연관성 X, 더욱 더 확실한 분리를 위해
	ser := _srvcheckRepo.NewESElasticsearchCheckHistoryRepository(_srvcheckConfig.App, esCli, json.MapWriter())
	ssr := _srvcheckRepo.NewESSwarmpitCheckHistoryRepository(_srvcheckConfig.App, esCli, json.MapWriter())
	scsr := _srvcheckRepo.NewESConsulCheckHistoryRepository(_srvcheckConfig.App, esCli, json.MapWriter())

	// srvcheck domain usecase
	seu := _srvcheckUcase.NewElasticsearchCheckUsecase(_srvcheckConfig.App, ser, _slk, _es)
	ssu := _srvcheckUcase.NewSwarmpitCheckUsecase(_srvcheckConfig.App, ssr, _slk, _dkr)
	scsu := _srvcheckUcase.NewConsulCheckUsecase(_srvcheckConfig.App, scsr, _slk, _csl, _dkr)

	// srvcheck domain delivery
	_srvcheckChanDelivery.NewElasticsearchCheckHandler(time.Tick(_srvcheckConfig.App.ESCheckDeliveryPingCycle()), seu)
	_srvcheckChanDelivery.NewSwarmpitCheckHandler(time.Tick(_srvcheckConfig.App.SwarmpitCheckDeliveryPingCycle()), ssu)

	runtime.Goexit()
}
