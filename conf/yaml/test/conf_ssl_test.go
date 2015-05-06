package test_yaml_conf

import (
	"testing"

	"github.com/seanpont/assert"
	"github.com/spikeekips/kasi"
	"github.com/spikeekips/kasi/conf"
)

func TestSSL(t *testing.T) {
	assert := assert.Assert(t)

	var yml string
	var setting *conf.CoreSetting
	var err error

	yml = loadFile("config_ssl.yml")
	setting, err = kasi.ParseConfig(yml)
	assert.Equal(err, nil)

	assert.Equal(setting.Services[0].SSL.Cert, "/secret/a.cert")
	assert.Equal(setting.Services[0].SSL.Key, "/secret/a.key")
	assert.Equal(setting.Services[0].SSL.Pem, "/secret/a.pem")

	yml = loadFile("config_ssl_blank.yml")
	setting, err = kasi.ParseConfig(yml)
	assert.NotNil(err)
}

func TestSSLInheritFromDefault(t *testing.T) {
	assert := assert.Assert(t)

	var yml string
	var setting *conf.CoreSetting
	var err error

	yml = loadFile("config_ssl_inherit_from_default_but_default_is_nil.yml")
	setting, err = kasi.ParseConfig(yml)
	assert.Nil(err)
	assert.Nil(setting.Services[0].SSL)

	yml = loadFile("config_ssl_inherit_from_default.yml")
	setting, err = kasi.ParseConfig(yml)
	assert.Nil(err)
	assert.Equal(setting.Services[0].SSL, setting.GetDefaultSetting().SSL)

	yml = loadFile("config_ssl_inherit_from_default_but_no.yml")
	setting, err = kasi.ParseConfig(yml)
	assert.Nil(err)
	assert.Nil(setting.Services[0].SSL)
}
