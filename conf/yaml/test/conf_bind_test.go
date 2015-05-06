package test_yaml_conf

import (
	"net"
	"testing"

	"github.com/seanpont/assert"
	"github.com/spikeekips/kasi"
)

func TestBind(t *testing.T) {
	assert := assert.Assert(t)

	yml := loadFile("config_multiple_bind.yml")

	setting, err := kasi.ParseConfig(yml)
	assert.Equal(err, nil)

	cases := map[string]*net.TCPAddr{
		"a.com": &net.TCPAddr{Port: 8000},
		"b.com": &net.TCPAddr{Port: 8000},
		"c.com": &net.TCPAddr{Port: 8000},
		"d.com": &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8000},
	}

	for _, service := range setting.Services {
		assert.Equal(cases[service.Hostnames[0]], service.Bind)
	}
}

func TestBindInheritFromDefault(t *testing.T) {
	assert := assert.Assert(t)

	yml := loadFile("conf_default_bind.yml")

	setting, err := kasi.ParseConfig(yml)
	assert.Equal(err, nil)

	cases := map[string]*net.TCPAddr{
		"a.com": setting.GetDefaultSetting().Bind,
		"b.com": &net.TCPAddr{Port: 90},
	}

	for _, service := range setting.Services {
		assert.Equal(cases[service.Hostnames[0]], service.Bind)
	}
}
