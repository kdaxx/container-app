package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fahedouch/go-logrotate"
	"github.com/kdaxx/container-app/app/api"
	"github.com/kdaxx/container-app/app/conf"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"strings"
)

// Formatter provides log format for logger
type Formatter struct {
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	// current log level
	logLevel := entry.Level
	time := entry.Time.Format("2006-01-02 15:04:05")
	fields := entry.Data
	var data []byte
	if len(fields) > 0 {
		var err error
		data, err = json.Marshal(fields)
		if err != nil {
			return nil, err
		}
	}
	var caller string
	if entry.HasCaller() {
		caller = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
	}

	msg := entry.Message

	return []byte(fmt.Sprintf("[%-7s][%s] %s %s%s\n",
		strings.ToUpper(logLevel.String()), time, caller, msg, string(data))), nil

}

type AppLogger struct {
	cfg    *conf.LogConfig
	rotate *logrotate.Logger
	logger *logrus.Logger
}

func (appLogger *AppLogger) PreInit() any {
	return func(cfg *conf.LogConfig, appConfig *conf.AppConfig,
		/*do init after config injection*/ c *ConfigInjector) error {
		appLogger.cfg = cfg
		err := appLogger.configLogger()
		if err != nil {
			return err
		}
		if appConfig.Mode == api.ReleaseMode {
			appLogger.setReleaseMode()
		} else {
			appLogger.setDebugMode()
		}

		return nil
	}
}

func (appLogger *AppLogger) configLogger() error {
	cfg := appLogger.cfg
	appLogger.logger = logrus.StandardLogger()
	appLogger.logger.SetFormatter(&Formatter{})

	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return errors.New(
			fmt.Sprintf("failed to parse level: %s, %v", cfg.Level, err))
	}
	appLogger.logger.SetLevel(level)

	// compatible system logger with info level
	log.SetOutput(appLogger.logger.WriterLevel(logrus.InfoLevel))

	// output both stdout and appLogger file
	appLogger.rotate = &logrotate.Logger{
		Filename:           cfg.Filepath,
		FilenameTimeFormat: cfg.Format,
		MaxBytes:           1024 * 1024 * cfg.MaxSize,
		MaxBackups:         cfg.MaxBackups,
		MaxAge:             cfg.MaxAge, //days
		Compress:           false,      // disabled by default
	}
	// output to stdout and appLogger file
	writer := io.MultiWriter(os.Stdout, appLogger.rotate)
	appLogger.logger.SetOutput(writer)
	return nil
}

func (appLogger *AppLogger) setReleaseMode() {
	appLogger.logger.SetReportCaller(false)
}

func (appLogger *AppLogger) setDebugMode() {
	appLogger.logger.SetReportCaller(true)
}

func (appLogger *AppLogger) BeforeAppStop(ctx context.Context) error {
	return appLogger.rotate.Close()
}

func NewAppLogger() *AppLogger {
	// default value here, for config inject failed
	return &AppLogger{}
}
