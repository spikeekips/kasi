package conf

import (
	"errors"
	"fmt"

	"code.google.com/p/go-uuid/uuid"

	"github.com/spikeekips/kasi/util"
)

type Setting interface {
	GetID() string
	String() string
}

// CoreSetting
type CoreSetting struct {
	ID             string
	Env            *EnvSetting
	Services       []*ServiceSetting
	DefaultSetting *ServiceSetting
	Middlewares    MiddlewaresSetting
}

func (setting *CoreSetting) GetID() string {
	if setting.ID == "" {
		setting.ID = uuid.NewUUID().String()
	}

	return setting.ID
}

func (setting *CoreSetting) String() string {
	return util.ToJson(setting)
}

func (setting *CoreSetting) GetDefaultSetting() *ServiceSetting {
	return setting.DefaultSetting
}

func (setting *CoreSetting) GetServicesByBind() map[string][]*ServiceSetting {
	servicesByBind := map[string][]*ServiceSetting{}
	for _, service := range setting.Services {
		bindString := service.Bind.String()
		if _, ok := servicesByBind[bindString]; !ok {
			servicesByBind[bindString] = []*ServiceSetting{}
		}
		servicesByBind[bindString] = append(servicesByBind[bindString], service)
	}

	return servicesByBind
}

func (setting *CoreSetting) GetServiceByID(id string) (*ServiceSetting, error) {
	for _, serviceSetting := range setting.Services {
		if serviceSetting.GetID() == id {
			return serviceSetting, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("not found serviceSetting by `%s`", id))
}

func (setting *CoreSetting) GetEndpointByID(id string) (*EndpointSetting, error) {
	for _, serviceSetting := range setting.Services {
		if endpointSetting, found := serviceSetting.EndpointByID[id]; found {
			return endpointSetting, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("not found endpoints by `%s`", id))
}
