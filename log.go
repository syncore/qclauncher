// QCLauncher by syncore <syncore@syncore.org> 2017
// https://github.com/syncore/qclauncher

package qclauncher

import (
	"fmt"

	log "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *log.SugaredLogger

func NewLogger() *log.SugaredLogger {
	logCfg := log.NewProductionConfig()
	logCfg.OutputPaths = []string{getLogFilePath()}
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
	return lgr.Sugar()
}

func setLogger() {
	if logger != nil {
		return
	}
	logger = NewLogger()
}
