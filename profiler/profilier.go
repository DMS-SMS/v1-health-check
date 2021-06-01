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
	return &defaultProfiler{
		awsSession: s,
		signalChan: make(chan os.Signal, 1),
		waitGroup:  wg,
		myCfg:      cfg,
	}
}

func (dp *defaultProfiler) StartProfiling() {
	now := time.Now()
	if now.Location().String() == time.UTC.String() {
		now = now.Add(time.Hour * 9)
	}
	nowDate := fmt.Sprintf("%4d-%02d-%02d", now.Year(), now.Month(), now.Day())
	nowTime := fmt.Sprintf("%02d:%02d:%02d", now.Hour(), now.Minute(), now.Second())

	// Ex) /usr/share/health-check/profile/v.1.0.0/2021-06-01/20:17:31/cpu.prof
	profPath := fmt.Sprintf("/usr/share/health-check/profile/v.%s/%s/%s", dp.myCfg.Version(), nowDate, nowTime)
	cpuProf := profPath + "/cpu.prof"
	memoryProf := profPath + "/memory.prof"
	blockProf := profPath + "/block.prof"

	if err := os.MkdirAll(profPath, os.ModePerm); err != nil {
		log.Fatalln(err)
	}
	cpuProfFile, err := os.Create(cpuProf)
	if err != nil {
		log.Fatalln(err)
	}
	memoryProfFile, err := os.Create(memoryProf)
	if err != nil {
		log.Fatalln(err)
	}
	blockProfFile, err := os.Create(blockProf)
	if err != nil {
		log.Fatalln(err)
	}

	// Ex) profiles/health-check/v.1.0.0/2021-06-01/20:17:31/cpu.prof
	profS3Path := fmt.Sprintf("profiling/health-check/v.%s/%s/%s", dp.myCfg.Version(), nowDate, nowTime)
	cpuProfS3 := profS3Path + "/cpu.prof"
	memoryProfS3 := profS3Path + "/memory.prof"
	blockProfS3 := profS3Path + "/block.prof"

	log.Println("start profiling cpu, memory, block")
	if err := pprof.StartCPUProfile(cpuProfFile); err != nil {
		log.Fatalln(err)
	}

	// 시작 당일이 끝날 때 까지 대기
	afterOneDay := now.AddDate(0, 0, 1)
	tomorrow := time.Date(afterOneDay.Year(), afterOneDay.Month(), afterOneDay.Day(), 0, 0, 0, 0, time.UTC)
	timeFinSig := time.Tick(tomorrow.Sub(now))

	select {
	case <-timeFinSig:
		log.Println("upload profiling result recorded on this day")
	case <-dp.signalChan:
		dp.waitGroup.Add(1)
		defer dp.waitGroup.Done()
		log.Println("upload profiling result as shutdown of process")
	}

	runtime.GC()
	_ = pprof.Lookup("heap").WriteTo(memoryProfFile, 1)
	_ = pprof.Lookup("block").WriteTo(blockProfFile, 1)
	pprof.StopCPUProfile()

	_ = cpuProfFile.Close()
	_ = memoryProfFile.Close()
	_ = blockProfFile.Close()

	cpuProfFile, _ = os.Open(cpuProf)
	memoryProfFile, _ = os.Open(memoryProf)
	blockProfFile, _ = os.Open(blockProf)

	svc := s3.New(dp.awsSession)
	if _, err := svc.PutObject(&s3.PutObjectInput{
		Body:   cpuProfFile,
		Bucket: aws.String(dp.myCfg.ProfilerS3Bucket()),
		Key:    aws.String(cpuProfS3),
	}); err != nil {
		log.Fatalf("unable to upload cpu profiling result to s3, err: %v\n", err)
	}
	if _, err := svc.PutObject(&s3.PutObjectInput{
		Body:   memoryProfFile,
		Bucket: aws.String(dp.myCfg.ProfilerS3Bucket()),
		Key:    aws.String(memoryProfS3),
	}); err != nil {
		log.Fatalf("unable to upload memory profiling result to s3, err: %v\n", err)
	}
	if _, err := svc.PutObject(&s3.PutObjectInput{
		Body:   blockProfFile,
		Bucket: aws.String(dp.myCfg.ProfilerS3Bucket()),
		Key:    aws.String(blockProfS3),
	}); err != nil {
		log.Fatalf("unable to upload block profiling result to s3, err: %v\n", err)
	}
}
