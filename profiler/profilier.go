package profiler

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

// defaultProfiler profile memory and heap usage & save to s3 per a day or when msg come from signal channel
type defaultProfiler struct {
	awsSession *session.Session
	signalChan chan os.Signal
	waitGroup  *sync.WaitGroup
	myCfg      defaultProfilerConfig
}

// defaultProfilerConfig is the config getter interface about default profiler
type defaultProfilerConfig interface {
	// ProfilerS3Bucket method returns string represent profiler s3 bucket
	ProfilerS3Bucket() string

	// Version method returns string represent version
	Version() string
}

func New(s *session.Session, wg *sync.WaitGroup, cfg defaultProfilerConfig) *defaultProfiler {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	return &defaultProfiler{
		awsSession: s,
		signalChan: make(chan os.Signal, 1),
		waitGroup:  wg,
		myCfg:      cfg,
	}
}
