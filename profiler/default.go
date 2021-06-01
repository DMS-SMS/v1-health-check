package profiler

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
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
