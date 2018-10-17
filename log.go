// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"fmt"
	"strings"

	"net/url"
	"os"

	log "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type qlogger struct {
	*log.SugaredLogger
}

var logger *qlogger

func NewLogger() *qlogger {
	logCfg := log.NewProductionConfig()
	winFileSink := func(u *url.URL) (log.Sink, error) {
		// https://github.com/uber-go/zap/issues/621
		return os.OpenFile(u.Path[1:], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	}
	if err := log.RegisterSink("winfile", winFileSink); err != nil {
		// logger in main will have already registered before overall logger
		if !strings.Contains(err.Error(), "already registered") {
			panic(fmt.Errorf("Could not register sink for structured logger, error: %s", err))
		}
	}
	logCfg.OutputPaths = []string{"winfile:///" + getLogFilePath()}
	logCfg.Encoding = "json"
	logCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if ConfDebug {
		logCfg.Level.SetLevel(zapcore.DebugLevel)
		logCfg.DisableStacktrace = false
		logCfg.DisableCaller = false
		logCfg.Development = true
	} else {
		logCfg.Level.SetLevel(zapcore.ErrorLevel)
		logCfg.DisableStacktrace = true
		logCfg.DisableCaller = true
		logCfg.Development = false
	}
	lgr, err := logCfg.Build()
	if err != nil {
		panic(fmt.Errorf("Could not initialize structured logger, error: %s", err))
	}
	defer lgr.Sync()
	return &qlogger{lgr.Sugar()}
}

func (l *qlogger) FatalUIw(msg string, keysAndValues ...interface{}) {
	l.Errorw(msg, keysAndValues)
	exitFromUI()
}

func setLogger() {
	if logger != nil {
		return
	}
	logger = NewLogger()
}
