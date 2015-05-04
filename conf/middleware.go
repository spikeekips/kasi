package kasi_conf

import (
	"errors"
	"fmt"

	"github.com/moraes/config"
	"github.com/spikeekips/kasi/util"
)

type MiddlewaresSetting map[string][]string

func (setting MiddlewaresSetting) GetID() string {
	return ""
}

func (setting MiddlewaresSetting) String() string {
	return kasi_util.ToJson(setting)
}

func parseMiddlewares(itemConfig *config.Config) (MiddlewaresSetting, error) {
	items, err := itemConfig.List("")
	if err != nil {
		return nil, err
	}

	middlewares := MiddlewaresSetting{}

	for i := 0; i < len(items); i++ {
		var name string
		{
			item, err := itemConfig.Map(fmt.Sprintf("%d", i))
			if err != nil {
				continue
			}
			// get name
			for k, _ := range item {
				name = k
			}
		}

		ms, err := itemConfig.List(fmt.Sprintf("%d.%s", i, name))
		if err != nil {
			return nil, err
		}

		middlewares[name] = MiddlewareSetting{}
		for j := 0; j < len(ms); j++ {
			m, err := itemConfig.String(fmt.Sprintf("%d.%s.%d", i, name, j))
			if err != nil {
				return nil, err
			}
			middlewares[name] = append(middlewares[name], m)
		}
	}

	return middlewares, nil
}

type MiddlewareSetting []string

func (setting MiddlewareSetting) String() string {
	return kasi_util.ToJson(setting)
}

func (setting MiddlewaresSetting) parseMiddlewareSetting(itemConfig *config.Config, defaultName string) (MiddlewareSetting, error) {
	enabled, err := itemConfig.Bool("middleware")
	if err == nil && !enabled {
		return nil, nil
	}

	middlewareName, err := itemConfig.String("middleware")
	if err != nil {
		if len(defaultName) < 1 {
			return nil, nil
		}
		middlewareName = defaultName
	}

	if middleware, found := setting[middlewareName]; found {
		return middleware, nil
	}

	return nil, errors.New(fmt.Sprintf("invalid middleware found, `%s`", middlewareName))
}
