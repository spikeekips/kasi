package kasi_t

import (
	"testing"

	"github.com/seanpont/assert"
	"github.com/spikeekips/kasi"
	"github.com/spikeekips/kasi/conf"
)

func TestEndpoints(t *testing.T) {
	assert := assert.Assert(t)

	var yml string
	var setting *kasi_conf.CoreSetting
	var err error

	yml = loadFile("config_simple_endpoint0.yml")

	setting, err = kasi.ParseConfig(yml)
	assert.Nil(err)
	assert.Equal(len(setting.Services[0].Endpoints), 1)

	yml = loadFile("config_simple_endpoint1.yml")

	setting, err = kasi.ParseConfig(yml)
	assert.Nil(err)
	assert.Equal(len(setting.Services[0].Endpoints), 2)
}

func TestDuplicatedEndpoints(t *testing.T) {
	assert := assert.Assert(t)

	var yml string
	var err error

	yml = loadFile("config_duplicated_endpoints.yml")

	_, err = kasi.ParseConfig(yml)
	assert.NotNil(err)
}

func TestEndpointsOpen(t *testing.T) {
	assert := assert.Assert(t)

	var yml string
	var setting *kasi_conf.CoreSetting
	var err error

	yml = loadFile("config_simple_endpoint0.yml")
	setting, err = kasi.ParseConfig(yml)
	assert.Nil(err)
	assert.Equal(len(setting.Services[0].Endpoints), 1)

	yml = loadFile("config_opened_endpoint.yml")
	setting, err = kasi.ParseConfig(yml)
	assert.Nil(err)
	assert.Equal(len(setting.Services[0].Endpoints), 1)

	yml = loadFile("config_closed_endpoint.yml")

	setting, err = kasi.ParseConfig(yml)
	assert.Nil(err)

	assert.False(setting.Services[0].Opened(), "")
	assert.Equal(len(setting.Services[0].Endpoints), 0)
}
