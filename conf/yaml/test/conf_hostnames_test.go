package test_yaml_conf

import (
	"sort"
	"testing"

	"github.com/seanpont/assert"
	"github.com/spikeekips/kasi"
)

func TestHostnamesWithList(t *testing.T) {
	assert := assert.Assert(t)

	yml := loadFile("config_hostnames_list_type.yml")
	setting, err := kasi.ParseConfig(yml)
	assert.Nil(err)

	assert.Equal(len(setting.Services), 1)
	assert.Equal(len(setting.Services[0].Hostnames), 2)

	for _, h := range []string{"my0.github.com", "my1.github.com"} {
		if sort.SearchStrings(setting.Services[0].Hostnames, h) == 2 {
			t.Error("failed to parse hostnames")
		}
	}
}

func TestHostnamesWithString(t *testing.T) {
	assert := assert.Assert(t)

	yml := loadFile("config_hostnames_string.yml")
	setting, err := kasi.ParseConfig(yml)
	assert.Nil(err)

	assert.Equal(len(setting.Services), 1)
	assert.Equal(len(setting.Services[0].Hostnames), 1)
	assert.Equal(setting.Services[0].Hostnames[0], "my0.github.com")
}

func TestHostnamesWithStrings(t *testing.T) {
	assert := assert.Assert(t)

	yml := loadFile("config_hostnames_strings.yml")
	setting, err := kasi.ParseConfig(yml)
	assert.Nil(err)

	assert.Equal(len(setting.Services), 1)
	assert.Equal(len(setting.Services[0].Hostnames), 2)

	for _, h := range []string{"my0.github.com", "my1.github.com"} {
		if sort.SearchStrings(setting.Services[0].Hostnames, h) == 2 {
			t.Error("failed to parse hostnames")
		}
	}
}

func TestHostnamesInheritFromDefault(t *testing.T) {
	assert := assert.Assert(t)

	yml := loadFile("config_hostnames_inherit_from_default.yml")
	setting, err := kasi.ParseConfig(yml)
	assert.Equal(err, nil)

	assert.Equal(len(setting.Services), 1)
	assert.Equal(len(setting.Services[0].Hostnames), 2)

	for _, h := range []string{"my0.github.com", "my1.github.com"} {
		if sort.SearchStrings(setting.Services[0].Hostnames, h) == 2 {
			t.Error("failed to parse hostnames")
		}
	}

	ymlOverride := loadFile("config_hostnames_override_default.yml")
	setting, err = kasi.ParseConfig(ymlOverride)
	assert.Nil(err)

	assert.Equal(len(setting.Services), 1)
	assert.Equal(len(setting.Services[0].Hostnames), 1)

	for _, h := range []string{"my2.github.com"} {
		if sort.SearchStrings(setting.Services[0].Hostnames, h) == 2 {
			t.Error("failed to parse hostnames")
		}
	}
}
