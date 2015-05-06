package kasi_conf

import (
	"os"
	"runtime"

	"github.com/moraes/config"
	"github.com/op/go-logging"
	"github.com/spikeekips/kasi/util"
)

type EnvSetting struct {
	Hostname            string
	GOOS                string
	MiddlewareDirectoty string
	LogLevel            logging.Level
}

func (setting *EnvSetting) String() string {
	return kasi_util.ToJson(setting)
}

func NewEnvSetting() *EnvSetting {
	setting := &EnvSetting{}

	setting.GOOS = runtime.GOOS
	setting.Hostname, _ = os.Hostname()

	return setting
}

func (setting *EnvSetting) parse(itemConfig *config.Config) error {
	_, err := itemConfig.Get("loglevel")
	if err != nil {
		return err
	}
	logLevelInput, err := itemConfig.String("loglevel")
	if err == nil {
		logLevel, err := kasi_util.GetLogLevel(logLevelInput)
		if err != nil {
			return err
		}
		setting.LogLevel = logLevel
	}
	return nil
}
