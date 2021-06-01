package profiler

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
	"os/signal"
	"syscall"
)

// defaultProfiler profile memory and heap usage & save to s3 per a day or when msg come from signal channel
type defaultProfiler struct {
	awsSession *session.Session
	signalChan chan<- os.Signal
	myCfg      defaultProfilerConfig
}

// defaultProfilerConfig is the config getter interface about default profiler
type defaultProfilerConfig interface {
	// ProfilerS3Bucket method returns string represent profiler s3 bucket
	ProfilerS3Bucket() string
}

func New(s *session.Session, cfg defaultProfilerConfig) *defaultProfiler {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	return &defaultProfiler{
		awsSession: s,
		signalChan: sig,
		myCfg:      cfg,
	}
}
