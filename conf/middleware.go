package conf

import "github.com/spikeekips/kasi/util"

type MiddlewaresSetting map[string][]string

func (setting MiddlewaresSetting) GetID() string {
	return ""
}

func (setting MiddlewaresSetting) String() string {
	return util.ToJson(setting)
}

type MiddlewareSetting []string

func (setting MiddlewareSetting) String() string {
	return util.ToJson(setting)
}
