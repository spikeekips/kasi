package test_yaml_conf

import (
	"testing"

	"github.com/seanpont/assert"
	"github.com/spikeekips/kasi"
	"github.com/spikeekips/kasi/conf"
)

func TestID(t *testing.T) {
	assert := assert.Assert(t)

	var yml string
	var setting *conf.CoreSetting
	var err error

	yml = loadFile("config_check_id.yml")
	setting, err = kasi.ParseConfig(yml)
	assert.Equal(err, nil)

	assert.Equal(setting.Services[0].GetID(), "this-is-service-id")
	assert.Equal(setting.Services[0].Endpoints[0].GetID(), "findme")
	assert.Equal(len(setting.Services[0].Endpoints[1].GetID()), 36)
	assert.Equal(setting.Services[0].Endpoints[1].GetID(), setting.Services[0].Endpoints[1].GetID())
}

func TestGetEndpointByID(t *testing.T) {
	assert := assert.Assert(t)

	var yml string
	var setting *conf.CoreSetting
	var endpointSetting *conf.EndpointSetting
	var err error

	yml = loadFile("config_check_id.yml")
	setting, err = kasi.ParseConfig(yml)
	assert.Equal(err, nil)

	endpointSetting, err = setting.GetEndpointByID("findme")
	assert.Equal(endpointSetting, setting.Services[0].Endpoints[0])

	endpointSettingToFind := setting.Services[0].Endpoints[1]
	endpointSetting, err = setting.GetEndpointByID(endpointSettingToFind.GetID())
	assert.Equal(endpointSettingToFind, endpointSetting)
}

func TestDuplicatedID(t *testing.T) {
	assert := assert.Assert(t)

	var yml string
	var err error

	yml = loadFile("config_duplicated_id.yml")
	_, err = kasi.ParseConfig(yml)
	assert.NotNil(err)
}
