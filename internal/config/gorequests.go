package config

import (
	"context"
	"os"

	"github.com/chyroc/gorequests"
	"github.com/sirupsen/logrus"
)

var ReqCli *gorequests.Factory

func init() {
	ReqCli = gorequests.NewFactory(gorequests.WithLogger(WithLogger()))
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
