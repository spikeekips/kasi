package kasi_conf

import (
	"os"
	"runtime"

	"github.com/moraes/config"
	"github.com/spikeekips/kasi/util"
)

type EnvSetting struct {
	Hostname            string
	GOOS                string
	MiddlewareDirectoty string
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

func (setting *EnvSetting) parse(itemconfig *config.Config) error {
	return nil
}
