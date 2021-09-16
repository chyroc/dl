package config

import (
	"context"
	"os"
	"time"

	"github.com/chyroc/gorequests"
	"github.com/sirupsen/logrus"
)

var ReqCli *gorequests.Factory

func init() {
	options := []gorequests.RequestOption{gorequests.WithLogger(WithLogger())}
	if os.Getenv("IN_CI") != "" {
		options = append(options, gorequests.WithTimeout(time.Second*10))
	}
	ReqCli = gorequests.NewFactory(options...)
}

func WithLogger() gorequests.Logger {
	logFile := WorkDir + "/log.log"

	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		panic(err)
	}
	logger := logrus.New()
	logger.SetOutput(f)

	return &gorequestsLogger{logger: logger}
}

type gorequestsLogger struct {
	logger *logrus.Logger
}

func (r *gorequestsLogger) Info(ctx context.Context, format string, v ...interface{}) {
	r.logger.Infof(format, v...)
}

func (r *gorequestsLogger) Error(ctx context.Context, format string, v ...interface{}) {
	r.logger.Errorf(format, v...)
}
