package kasi_conf

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"code.google.com/p/go-uuid/uuid"

	"github.com/moraes/config"
	"github.com/spikeekips/kasi/util"
)

type Setting interface {
	GetID() string
	String() string
}

// CoreSetting
type CoreSetting struct {
	ID             string
	conf           *config.Config
	Env            *EnvSetting
	Services       []*ServiceSetting
	defaultSetting *ServiceSetting
	Middlewares    MiddlewaresSetting
}

func (setting *CoreSetting) GetID() string {
	if setting.ID == "" {
		setting.ID = uuid.NewUUID().String()
	}

	return setting.ID
}

func (setting *CoreSetting) String() string {
	return kasi_util.ToJson(setting)
}

func (setting *CoreSetting) GetDefaultSetting() *ServiceSetting {
	return setting.defaultSetting
}

func (setting *CoreSetting) parseYamlDocument() error {
	items, err := setting.conf.List("")
	if err != nil {
		return err
	}

	if len(items) < 1 {
		return errors.New("No configuration found")
	}

	// Env
	env := NewEnvSetting()
	for i := 0; i < len(items); i++ {
		itemConfig, err := setting.conf.Get(fmt.Sprintf("%d.env", i))
		if err != nil {
			continue
		}

		err = env.parse(itemConfig)
		if err != nil {
			return err
		}
		break
	}
	setting.Env = env

	// Middlewares
	var middlewares MiddlewaresSetting
	for i := 0; i < len(items); i++ {
		itemConfig, err := setting.conf.Get(fmt.Sprintf("%d.middlewares", i))
		if err != nil {
			continue
		}

		middlewares, err = parseMiddlewares(itemConfig)
		if err != nil {
			return err
		}
		break
	}
	setting.Middlewares = middlewares

	// setting.defaultSetting
	setting.defaultSetting = &ServiceSetting{}
	for i := 0; i < len(items); i++ {
		itemConfig, err := setting.conf.Get(fmt.Sprintf("%d.default", i))
		if err != nil {
			continue
		}

		serviceSetting, err := setting.parseDefaultSetting(itemConfig)
		if err != nil {
			return err
		}
		setting.defaultSetting = serviceSetting
		break
	}

	// setting.ServiceSetting
	for i := 0; i < len(items); i++ {
		itemConfig, err := setting.conf.Get(fmt.Sprintf("%d.service", i))
		if err != nil {
			continue
		}

		serviceSetting, err := setting.parseServiceSetting(itemConfig)
		if err != nil {
			return err
		}
		setting.Services = append(setting.Services, serviceSetting)
	}

	if len(setting.Services) < 1 {
		return errors.New("no service found")
	}

	// check id duplication
	ids := []string{}
	for _, serviceSetting := range setting.Services {
		id := serviceSetting.GetID()
		if kasi_util.InArray(ids, id) {
			return errors.New(fmt.Sprintf("found duplicated id, `%s`, %v", id, ids))
		}

		ids = append(ids, id)

		for _, endpointSetting := range serviceSetting.Endpoints {
			id = endpointSetting.GetID()
			if kasi_util.InArray(ids, id) {
				return errors.New(fmt.Sprintf("found duplicated id, `%s`, %v", id, ids))
			}

			ids = append(ids, id)
		}
	}

	return nil
}

func (setting *CoreSetting) parseSSLSetting(itemConfig *config.Config) (*SSLSetting, error) {
	_, err := itemConfig.Map("ssl")
	if err != nil {
		return nil, nil
	}

	sslCert, _ := itemConfig.String("ssl.cert")
	sslKey, _ := itemConfig.String("ssl.key")
	sslPem, _ := itemConfig.String("ssl.pem")

	if len(kasi_util.Trim(sslCert)) < 1 && len(kasi_util.Trim(sslKey)) < 1 && len(kasi_util.Trim(sslKey)) < 1 {
		return nil, errors.New("invalid ssl")
	}

	return &SSLSetting{Cert: sslCert, Key: sslKey, Pem: sslPem}, nil
}

func (setting *CoreSetting) parseBaseServiceSetting(itemConfig *config.Config) (*ServiceSetting, error) {
	serviceSetting := ServiceSetting{}

	var hostnameStrings []string

	rawList, err := itemConfig.List("hostname")
	if err == nil {
		for _, s := range rawList {
			hostnameStrings = append(hostnameStrings, s.(string))
		}
	} else {
		rawString, err := itemConfig.String("hostname")
		if err == nil {
			for _, s := range strings.Split(rawString, ",") {
				striped := kasi_util.Trim(s)
				if len(striped) < 1 {
					continue
				}

				hostnameStrings = append(hostnameStrings, striped)
			}
		}
	}

	serviceSetting.Hostnames = hostnameStrings

	// bind
	// TODO support all available `net.Addr`
	rawBind, err := itemConfig.String("bind")
	if err == nil {
		bindAddr, err := splitHostPort(rawBind)
		if err != nil {
			return nil, err
		}
		serviceSetting.Bind = bindAddr
	}

	// source
	sources, err := parseSource(itemConfig)
	if err == nil {
		serviceSetting.Sources = sources
	}

	// timeout
	var timeout time.Duration
	rawTimeout, err := itemConfig.String("timeout")
	if err != nil {
		timeout = time.Nanosecond * -1
	} else {
		timeout, err = parseTimeUnit(rawTimeout)
		if err != nil {
			return nil, err
		}
	}
	serviceSetting.Timeout = timeout

	// ssl
	ssl, err := setting.parseSSLSetting(itemConfig)
	if err != nil {
		return nil, err
	}

	serviceSetting.SSL = ssl

	// middleware
	middlewareName, err := itemConfig.String("middleware")
	if err == nil {
		if _, found := setting.Middlewares[middlewareName]; !found {
			return nil, errors.New(fmt.Sprintf("invalid middleware found, `%s`", middlewareName))
		}
		serviceSetting.Middleware = middlewareName
	}

	return &serviceSetting, nil
}

func (setting *CoreSetting) parseDefaultSetting(itemConfig *config.Config) (*ServiceSetting, error) {
	serviceSetting, err := setting.parseBaseServiceSetting(itemConfig)
	if err != nil {
		return nil, err
	}

	if serviceSetting.Timeout < 0 {
		serviceSetting.Timeout = time.Nanosecond * 0
	}

	// the sources in default must be valid url
	if !validateSources(serviceSetting.Sources) {
		return nil, errors.New("the sources in the defaults, must be valid url")
	}

	return serviceSetting, nil
}

func (setting *CoreSetting) parseServiceSetting(itemConfig *config.Config) (*ServiceSetting, error) {
	serviceSetting, err := setting.parseBaseServiceSetting(itemConfig)
	if err != nil {
		return nil, err
	}

	// id
	id, err := itemConfig.String("id")
	if err == nil {
		serviceSetting.ID = id
	}

	// check hostnames must be exists.
	if len(serviceSetting.Hostnames) < 1 {
		if len(setting.defaultSetting.Hostnames) < 1 {
			return nil, errors.New("failed to set hostname")
		}
		serviceSetting.Hostnames = setting.defaultSetting.Hostnames
	}

	if serviceSetting.Bind == nil {
		if setting.defaultSetting.Bind == nil {
			return nil, errors.New("failed to set bind")
		}
		serviceSetting.Bind = setting.defaultSetting.Bind
	}

	var sources []string
	if serviceSetting.Sources != nil {
		sources = MergeURLs(setting.defaultSetting.Sources, serviceSetting.Sources)
	} else {
		if len(setting.defaultSetting.Sources) == 0 {
			return nil, errors.New("failed to set sources")
		}
		sources = setting.defaultSetting.Sources
	}

	if len(sources) < 1 || !validateSources(sources) {
		return nil, errors.New("sources not found")
	}

	serviceSetting.Sources = sources

	if serviceSetting.Timeout < 0 {
		serviceSetting.Timeout = setting.defaultSetting.Timeout
	}

	if serviceSetting.SSL == nil {
		sslEnabled, err := itemConfig.Bool("ssl")
		if err == nil {
			if !sslEnabled {
				serviceSetting.SSL = nil
			} else {
				if setting.defaultSetting.SSL != nil {
					serviceSetting.SSL = setting.defaultSetting.SSL
				}
			}
		}
	}

	// EndpointSetting
	endpoints, err := itemConfig.List("endpoints")
	if err != nil {
		return nil, err
	}

	if len(endpoints) < 1 {
		return nil, errors.New("No endpoints found")
	}

	var endpointsSttingsParsed []*EndpointSetting
	for i := 0; i < len(endpoints); i++ {
		endpointConfig, err := itemConfig.Get(fmt.Sprintf("endpoints.%d.endpoint", i))
		if err != nil {
			return nil, err
		}

		endpointSetting, err := setting.parseEndpointSetting(serviceSetting, endpointConfig)
		if err != nil {
			return nil, err
		}
		if endpointSetting == nil || !endpointSetting.Opened() {
			continue
		}

		endpointsSttingsParsed = append(endpointsSttingsParsed, endpointSetting)
	}
	endpointSettings := EndpointSettings(endpointsSttingsParsed)

	if len(kasi_util.GetUnique(endpointSettings)) != len(endpointsSttingsParsed) {
		return nil, errors.New("duplicated endpoints found")
	}

	serviceSetting.Endpoints = endpointSettings

	serviceSetting.EndpointByID = map[string]*EndpointSetting{}
	for _, endpointSetting := range serviceSetting.Endpoints {
		serviceSetting.EndpointByID[endpointSetting.GetID()] = endpointSetting
	}

	return serviceSetting, nil
}

func (setting *CoreSetting) parseEndpointSetting(
	serviceSetting *ServiceSetting,
	itemConfig *config.Config,
) (*EndpointSetting, error) {
	endpointSetting := &EndpointSetting{}

	// open
	_, err := itemConfig.Get("open")
	if err != nil {
		endpointSetting.open = true
	} else {
		opened, err := itemConfig.Bool("open")
		if err != nil {
			return nil, err
		}
		endpointSetting.open = opened
	}

	// id
	id, err := itemConfig.String("id")
	if err == nil {
		endpointSetting.ID = id
	}

	// expose
	rawExpose, err := itemConfig.String("expose")
	if err != nil {
		return nil, err
	}
	if len(kasi_util.Trim(rawExpose)) < 1 {
		return nil, errors.New("empty expose")
	}
	endpointSetting.Expose = kasi_util.Trim(rawExpose)
	endpointSetting.exposeRegexp, err = regexp.Compile(endpointSetting.Expose)
	if err != nil {
		return nil, err
	}

	// source
	parsedSources, err := parseSource(itemConfig)
	if err != nil {
		return nil, err
	}
	sources := MergeURLs(serviceSetting.Sources, parsedSources)
	if !validateSources(sources) {
		return nil, errors.New("the sources in the defaults, must be valid url")
	}
	endpointSetting.Sources = sources

	// middleware
	middleware, err := setting.Middlewares.parseMiddlewareSetting(itemConfig, serviceSetting.Middleware)
	if err != nil {
		return nil, err
	}
	endpointSetting.Middleware = middleware

	return endpointSetting, nil
}

func ParseConfig(conf *config.Config) (*CoreSetting, error) {
	setting := &CoreSetting{conf: conf}
	err := setting.parseYamlDocument()
	if err != nil {
		return nil, err
	}

	return setting, nil
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
