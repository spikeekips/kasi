package kasi_t

import (
	"testing"

	"github.com/seanpont/assert"
	"github.com/spikeekips/kasi"
)

func TestParseYaml(t *testing.T) {
	assert := assert.Assert(t)

	var yml string = `
	`
	_, err := kasi.ParseConfig(yml)
	assert.NotNil(err)

	yml = loadFile("config_bad_formatted_without_doc_declaration.yml")
	_, err = kasi.ParseConfig(yml)
	assert.NotNil(err)

	yml = loadFile("config_well_formatted0.yml")
	_, err = kasi.ParseConfig(yml)
	assert.Equal(err, nil)
}

func TestParseSimple(t *testing.T) {
	assert := assert.Assert(t)

	yml := loadFile("config_well_formatted1.yml")
	setting, err := kasi.ParseConfig(yml)
	assert.Nil(err)
	assert.Equal(len(setting.Services), 1)
	assert.NotNil(setting.GetDefaultSetting())
}
