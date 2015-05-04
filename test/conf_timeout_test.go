package kasi_t

import (
	"testing"
	"time"

	"github.com/seanpont/assert"
	"github.com/spikeekips/kasi"
)

func TestTimeout(t *testing.T) {
	assert := assert.Assert(t)

	yml := loadFile("config_timeout.yml")
	setting, err := kasi.ParseConfig(yml)
	assert.Equal(err, nil)

	cases := map[string]time.Duration{
		"a.com": time.Second * 9,
		"b.com": time.Second * 0,
		"c.com": setting.GetDefaultSetting().Timeout,
		"d.com": setting.GetDefaultSetting().Timeout,
	}

	for _, service := range setting.Services {
		assert.Equal(service.Timeout, cases[service.Hostnames[0]])
	}
}
