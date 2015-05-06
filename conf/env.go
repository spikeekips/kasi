package conf

import (
	"os"
	"runtime"

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
	return util.ToJson(setting)
}

func NewEnvSetting() *EnvSetting {
	setting := &EnvSetting{}

	setting.GOOS = runtime.GOOS
	setting.Hostname, _ = os.Hostname()

	return setting
}
