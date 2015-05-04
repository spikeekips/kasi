package kasi_t

import (
	"testing"

	"github.com/seanpont/assert"
	"github.com/spikeekips/kasi"
	"github.com/spikeekips/kasi/conf"
)

func TestSource(t *testing.T) {
	assert := assert.Assert(t)

	yml := loadFile("config_source.yml")
	setting, err := kasi.ParseConfig(yml)
	assert.Equal(err, nil)

	cases := map[string]string{
		"a.com": "https://a.com/api/v1",
		"b.com": "https://b.com/api/v1",
	}

	for _, service := range setting.Services {
		assert.Equal(cases[service.Hostnames[0]], service.Sources[0])
	}
}

func TestSourceInvalidURL(t *testing.T) {
	assert := assert.Assert(t)

	var yml string

	yml = loadFile("config_source_invalid_url0.yml")
	_, err := kasi.ParseConfig(yml)
	assert.NotNil(err)

	yml = loadFile("config_source_invalid_url1.yml")
	_, err = kasi.ParseConfig(yml)
	assert.NotNil(err)
}

func TestSourceList(t *testing.T) {
	assert := assert.Assert(t)

	var yml string
	var setting *kasi_conf.CoreSetting
	var err error

	yml = loadFile("config_source_list_type.yml")
	setting, err = kasi.ParseConfig(yml)
	assert.Nil(err)
	assert.Equal(len(setting.Services[0].Sources), 2)

	merged := kasi_conf.MergeURLs(setting.Services[0].Sources, []string{"/b"})

	assert.Equal(setting.Services[0].Endpoints[0].Sources, merged)
}
