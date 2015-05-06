package kasi

import (
	"github.com/spikeekips/kasi/conf"
	"github.com/spikeekips/kasi/conf/yaml"
)

// parse functions
func ParseConfig(s string) (*conf.CoreSetting, error) {
	setting, err := conf_yaml.ParseConfig(s)
	if err != nil {
		return nil, err
	}

	return setting, nil
}
