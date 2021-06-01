// Create package in v.1.0.0 (main)
// Most package import and dependency object generation, injection occurs in this package

package main

import (
	// import Go SDK package
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	// import external package
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/docker/docker/client"
	es "github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	// import app config & various agent package
	"github.com/DMS-SMS/v1-health-check/app/config"
	"github.com/DMS-SMS/v1-health-check/consul"
	"github.com/DMS-SMS/v1-health-check/docker"
	"github.com/DMS-SMS/v1-health-check/elasticsearch"
	"github.com/DMS-SMS/v1-health-check/grpc"
	"github.com/DMS-SMS/v1-health-check/json"
	"github.com/DMS-SMS/v1-health-check/profiler"
	"github.com/DMS-SMS/v1-health-check/slack"
	"github.com/DMS-SMS/v1-health-check/system"

	// import system check domain package
	_syscheckConfig "github.com/DMS-SMS/v1-health-check/syscheck/config"
	_syscheckChanDelivery "github.com/DMS-SMS/v1-health-check/syscheck/delivery/channel"
	_syscheckHttpDelivery "github.com/DMS-SMS/v1-health-check/syscheck/delivery/http"
	_syscheckRepo "github.com/DMS-SMS/v1-health-check/syscheck/repository/elasticsearch"
	_syscheckUcase "github.com/DMS-SMS/v1-health-check/syscheck/usecase"

	// import service check domain package
	_srvcheckConfig "github.com/DMS-SMS/v1-health-check/srvcheck/config"
	_srvcheckChanDelivery "github.com/DMS-SMS/v1-health-check/srvcheck/delivery/channel"
	_srvcheckHttpDelivery "github.com/DMS-SMS/v1-health-check/srvcheck/delivery/http"
	_srvcheckRepo "github.com/DMS-SMS/v1-health-check/srvcheck/repository/elasticsearch"
	_srvcheckUcase "github.com/DMS-SMS/v1-health-check/srvcheck/usecase"
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
	// add elasticsearch API connection
	esCli, err := es.NewClient(es.Config{
		Addresses: []string{config.App.ESAddress()},
	})
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create elasticsearch client"))
	}

	// add docker engine API connection
	dkrCli, err := client.NewClientWithOpts(
		client.WithVersion(config.App.DockerCliVer()),
		client.WithTimeout(config.App.DockerCliTimeout()),
	)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create docker client"))
	}

	// add consul API connection
	cslCfg := api.DefaultConfig()
	cslCfg.Address = config.App.ConsulAddress()
	cslCli, err := api.NewClient(cslCfg)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create consul client"))
	}

	// add aws session connection
	awsSess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.App.AWSRegion()),
		Credentials: credentials.NewStaticCredentials(config.App.AWSAccountID(), config.App.AWSAccountKey(), ""),
	})
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create aws session"))
	}

	// define ctx having WaitGroup in value & which is type of cancelCtx
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), "WaitGroup", wg))

	// start profiling with goroutine until stop process
	prof := profiler.New(awsSess, wg, config.App)
	go func(profileFunc func()) {
		for {
			if ctx.Err() == context.Canceled {
				break
			}
			profileFunc()
		}
	}(prof.StartProfiling)

	// add docker, system, slack, elasticsearch agent
	_dkr := docker.NewAgent(dkrCli)
	_sys := system.NewAgent(dkrCli)
	_slk := slack.NewAgent(config.App.SlackAPIToken(), config.App.SlackChatChannel())
	_es := elasticsearch.NewAgent(esCli)
	_csl := consul.NewAgent(cslCli)
	_rpc := grpc.NewGRPCAgent()

	// about syscheck domain
	// syscheck domain repository
	sdr := _syscheckRepo.NewESDiskCheckHistoryRepository(_syscheckConfig.App, esCli, json.MapWriter())
	scr := _syscheckRepo.NewESCPUCheckHistoryRepository(_syscheckConfig.App, esCli, json.MapWriter())
	smr := _syscheckRepo.NewESMemoryCheckHistoryRepository(_syscheckConfig.App, esCli, json.MapWriter())

	// syscheck domain usecase
	sdu := _syscheckUcase.NewDiskCheckUsecase(_syscheckConfig.App, sdr, _slk, _sys)
	scu := _syscheckUcase.NewCPUCheckUsecase(_syscheckConfig.App, scr, _slk, _sys, _dkr)
	smu := _syscheckUcase.NewMemoryCheckUsecase(_syscheckConfig.App, smr, _slk, _sys, _dkr)

	// syscheck domain delivery
	_syscheckChanDelivery.SetGlobalContext(ctx)
	_syscheckChanDelivery.NewDiskCheckHandler(time.Tick(_syscheckConfig.App.DiskCheckDeliveryPingCycle()), sdu)
	_syscheckChanDelivery.NewCPUCheckHandler(time.Tick(_syscheckConfig.App.CPUCheckDeliveryPingCycle()), scu)
	_syscheckChanDelivery.NewMemoryCheckHandler(time.Tick(_syscheckConfig.App.MemoryCheckDeliveryPingCycle()), smu)

	// about srvcheck domain
	// srvcheck domain repository
	ser := _srvcheckRepo.NewESElasticsearchCheckHistoryRepository(_srvcheckConfig.App, esCli, json.MapWriter())
	ssr := _srvcheckRepo.NewESSwarmpitCheckHistoryRepository(_srvcheckConfig.App, esCli, json.MapWriter())
	scsr := _srvcheckRepo.NewESConsulCheckHistoryRepository(_srvcheckConfig.App, esCli, json.MapWriter())

	// srvcheck domain usecase
	seu := _srvcheckUcase.NewElasticsearchCheckUsecase(_srvcheckConfig.App, ser, _slk, _es)
	ssu := _srvcheckUcase.NewSwarmpitCheckUsecase(_srvcheckConfig.App, ssr, _slk, _dkr)
	scsu := _srvcheckUcase.NewConsulCheckUsecase(_srvcheckConfig.App, scsr, _slk, _csl, _rpc, _dkr)

	// srvcheck domain delivery
	_srvcheckChanDelivery.SetGlobalContext(ctx)
	_srvcheckChanDelivery.NewElasticsearchCheckHandler(time.Tick(_srvcheckConfig.App.ESCheckDeliveryPingCycle()), seu)
	_srvcheckChanDelivery.NewSwarmpitCheckHandler(time.Tick(_srvcheckConfig.App.SwarmpitCheckDeliveryPingCycle()), ssu)
	_srvcheckChanDelivery.NewConsulCheckHandler(time.Tick(_srvcheckConfig.App.ConsulCheckDeliveryPingCycle()), scsu)

	// expose usecase method to HTTP API
	r := gin.Default()
	_syscheckHttpDelivery.NewSyscheckHandler(r, sdu, scu, smu)
	_srvcheckHttpDelivery.NewSrvcheckHandler(r, scsu, seu, ssu)

	gin.SetMode(gin.ReleaseMode)
	go func() { _ = r.Run(":8888") }()

	// handle signal to graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	s := <-sigs
	log.Printf("SIGNAL TO STOP PROCESS WAS NOTIFIED!, SIGNAL: %s", s)

	cancel()
	log.Println("CANCEL DELIVERY CONTEXT & WAIT TO ALL HANDLING GROUP DONE!")

	prof.StopProfiling(s)
	log.Println("STOP PROFILING & SAVE PROFILE RESULT TO AWS S3!")

	wg.Wait()
	log.Println("ALL HANDLING GROUP WAS DONE! SUCCEED TO GRACEFUL SHUTDOWN.")
}
