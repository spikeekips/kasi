package kasi_t

import (
	"testing"

	"github.com/seanpont/assert"
	"github.com/spikeekips/kasi"
)

func TestEnv(t *testing.T) {
	assert := assert.Assert(t)

	yml := loadFile("config_env.yml")
	setting, err := kasi.ParseConfig(yml)
	assert.Equal(err, nil)

	assert.NotEqual(setting.Env.Hostname, "")
	assert.NotEqual(setting.Env.GOOS, "")
}
