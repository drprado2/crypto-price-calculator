package applog

import (
	"context"
	"crypto-price-calculator/internal/configs"
	log "github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"os"
	"time"
)

func Setup() {
	//envs := configs.Get()
	//if envs.ServerEnvironment == configs.DeveloperEnvironment || envs.ServerEnvironment == configs.TestEnvironment {
	//} else {
	//	log.SetFormatter(&log.JSONFormatter{})
	//	log.SetLevel(log.InfoLevel)
	//}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:             true,
		ForceColors:               true,
		EnvironmentOverrideColors: true,
		PadLevelText:              true,
		TimestampFormat:           "2006-01-02 15:04:05",
	})
	log.SetLevel(log.DebugLevel)

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   "eng/logs/console.log",
		MaxSize:    50, // megabytes
		MaxBackups: 3,  // amouts
		MaxAge:     28, //days
		Level:      log.InfoLevel,
		Formatter: &log.JSONFormatter{
			TimestampFormat: time.RFC822,
		},
	})
	if err != nil {
		panic(err)
	}

	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)
	log.AddHook(rotateFileHook)
}

func Logger(ctx context.Context) *log.Entry {
	contextLogger := log.WithFields(getDefaultLogFields(ctx))
	return contextLogger
}

func getDefaultLogFields(ctx context.Context) log.Fields {
	envs := configs.Get()
	cid := ctx.Value("cid")
	if cid == nil {
		cid = "empty"
	}
	httpMethod := ctx.Value("httpMethod")
	httpPath := ctx.Value("httpPath")

	fields := log.Fields{
		"cid":     cid,
		"version": envs.SystemVersion,
		"app":     envs.AppName,
	}
	if httpMethod != nil {
		fields["httpMethod"] = httpMethod
	}
	if httpPath != nil {
		fields["httpPath"] = httpPath
	}
	return fields
}
