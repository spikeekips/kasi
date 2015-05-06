package conf_yaml

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/moraes/config"
	"github.com/op/go-logging"
	"github.com/spikeekips/kasi/conf"
	"github.com/spikeekips/kasi/util"
)

var RE_STRIP_TABED_LINE *regexp.Regexp = func() *regexp.Regexp { r, _ := regexp.Compile("^[\\s\\t]*$"); return r }()
var RE_STRIP_PREFIXED_TABS *regexp.Regexp = func() *regexp.Regexp { r, _ := regexp.Compile("^[\\t][\\t]*"); return r }()
var RE_STRIP_BLANKS *regexp.Regexp = func() *regexp.Regexp { r, _ := regexp.Compile("[\\s\\t]*$"); return r }()

func ParseConfig(yaml string) (*conf.CoreSetting, error) {
	// clean up the yaml string
	var cleanedYamlLine []string
	var cleanedYaml string
	var err error

	for _, s := range strings.Split(yaml, "\n") {
		if RE_STRIP_TABED_LINE.Match([]byte(s)) {
			cleanedYamlLine = append(cleanedYamlLine, "")
			continue
		}

		striped := RE_STRIP_BLANKS.ReplaceAllString(s, "")
		cleanedYamlLine = append(cleanedYamlLine, striped)
	}
	cleanedYaml = strings.Join(cleanedYamlLine, "\n")

	configParsed, err := config.ParseYaml(cleanedYaml)
	if err != nil {
		return nil, err
	}

	setting := &conf.CoreSetting{}
	err = parseYamlDocument(setting, configParsed)
	if err != nil {
		return nil, err
	}

	return setting, nil
}

func parseYamlDocument(setting *conf.CoreSetting, configParsed *config.Config) error {
	items, err := configParsed.List("")
	if err != nil {
		return err
	}

	if len(items) < 1 {
		return errors.New("No configuration found")
	}

	// Env
	envSetting := conf.NewEnvSetting()
	for i := 0; i < len(items); i++ {
		itemConfig, err := configParsed.Get(fmt.Sprintf("%d.env", i))
		if err != nil {
			continue
		}

		envSettingParsed, err := parseEnvSetting(setting, itemConfig)
		if err != nil {
			return err
		}
		envSetting = envSettingParsed
		break
	}
	setting.Env = envSetting

	// Middlewares
	var middlewares conf.MiddlewaresSetting
	for i := 0; i < len(items); i++ {
		itemConfig, err := configParsed.Get(fmt.Sprintf("%d.middlewares", i))
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
	setting.DefaultSetting = &conf.ServiceSetting{}
	for i := 0; i < len(items); i++ {
		itemConfig, err := configParsed.Get(fmt.Sprintf("%d.default", i))
		if err != nil {
			continue
		}

		serviceSetting, err := parseDefaultSetting(setting, itemConfig)
		if err != nil {
			return err
		}
		setting.DefaultSetting = serviceSetting
		break
	}

	// setting.ServiceSetting
	for i := 0; i < len(items); i++ {
		itemConfig, err := configParsed.Get(fmt.Sprintf("%d.service", i))
		if err != nil {
			continue
		}

		serviceSetting, err := parseServiceSetting(setting, itemConfig)
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
		if util.InArray(ids, id) {
			return errors.New(fmt.Sprintf("found duplicated id, `%s`, %v", id, ids))
		}

		ids = append(ids, id)

		for _, endpointSetting := range serviceSetting.Endpoints {
			id = endpointSetting.GetID()
			if util.InArray(ids, id) {
				return errors.New(fmt.Sprintf("found duplicated id, `%s`, %v", id, ids))
			}

			ids = append(ids, id)
		}
	}

	return nil
}

func parseSSLSetting(setting *conf.CoreSetting, itemConfig *config.Config) (*conf.SSLSetting, error) {
	_, err := itemConfig.Map("ssl")
	if err != nil {
		return nil, nil
	}

	sslCert, _ := itemConfig.String("ssl.cert")
	sslKey, _ := itemConfig.String("ssl.key")
	sslPem, _ := itemConfig.String("ssl.pem")

	if len(util.Trim(sslCert)) < 1 && len(util.Trim(sslKey)) < 1 && len(util.Trim(sslKey)) < 1 {
		return nil, errors.New("invalid ssl")
	}

	return &conf.SSLSetting{Cert: sslCert, Key: sslKey, Pem: sslPem}, nil
}

func parseBaseServiceSetting(setting *conf.CoreSetting, itemConfig *config.Config) (*conf.ServiceSetting, error) {
	serviceSetting := conf.ServiceSetting{}

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
				striped := util.Trim(s)
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
		bindAddr, err := conf.SplitHostPort(rawBind)
		if err != nil {
			return nil, err
		}
		serviceSetting.Bind = bindAddr
	}

	// source
	sources, err := conf.ParseSource(itemConfig)
	if err == nil {
		serviceSetting.Sources = sources
	}

	// timeout
	var timeout time.Duration
	rawTimeout, err := itemConfig.String("timeout")
	if err != nil {
		timeout = time.Nanosecond * -1
	} else {
		timeout, err = conf.ParseTimeUnit(rawTimeout)
		if err != nil {
			return nil, err
		}
	}
	serviceSetting.Timeout = timeout

	// ssl
	ssl, err := parseSSLSetting(setting, itemConfig)
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

func parseDefaultSetting(setting *conf.CoreSetting, itemConfig *config.Config) (*conf.ServiceSetting, error) {
	serviceSetting, err := parseBaseServiceSetting(setting, itemConfig)
	if err != nil {
		return nil, err
	}

	if serviceSetting.Timeout < 0 {
		serviceSetting.Timeout = time.Nanosecond * 0
	}

	// the sources in default must be valid url
	if !conf.ValidateSources(serviceSetting.Sources) {
		return nil, errors.New("the sources in the defaults, must be valid url")
	}

	return serviceSetting, nil
}

func parseServiceSetting(setting *conf.CoreSetting, itemConfig *config.Config) (*conf.ServiceSetting, error) {
	serviceSetting, err := parseBaseServiceSetting(setting, itemConfig)
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
		if len(setting.DefaultSetting.Hostnames) < 1 {
			return nil, errors.New("failed to set hostname")
		}
		serviceSetting.Hostnames = setting.DefaultSetting.Hostnames
	}

	if serviceSetting.Bind == nil {
		if setting.DefaultSetting.Bind == nil {
			return nil, errors.New("failed to set bind")
		}
		serviceSetting.Bind = setting.DefaultSetting.Bind
	}

	var sources []string
	if serviceSetting.Sources != nil {
		sources = conf.MergeURLs(setting.DefaultSetting.Sources, serviceSetting.Sources)
	} else {
		if len(setting.DefaultSetting.Sources) == 0 {
			return nil, errors.New("failed to set sources")
		}
		sources = setting.DefaultSetting.Sources
	}

	if len(sources) < 1 || !conf.ValidateSources(sources) {
		return nil, errors.New("sources not found")
	}

	serviceSetting.Sources = sources

	if serviceSetting.Timeout < 0 {
		serviceSetting.Timeout = setting.DefaultSetting.Timeout
	}

	if serviceSetting.SSL == nil {
		sslEnabled, err := itemConfig.Bool("ssl")
		if err == nil {
			if !sslEnabled {
				serviceSetting.SSL = nil
			} else {
				if setting.DefaultSetting.SSL != nil {
					serviceSetting.SSL = setting.DefaultSetting.SSL
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

	var endpointsSttingsParsed []*conf.EndpointSetting
	for i := 0; i < len(endpoints); i++ {
		endpointConfig, err := itemConfig.Get(fmt.Sprintf("endpoints.%d.endpoint", i))
		if err != nil {
			return nil, err
		}

		endpointSetting, err := parseEndpointSetting(setting, serviceSetting, endpointConfig)
		if err != nil {
			return nil, err
		}
		if endpointSetting == nil || !endpointSetting.Open {
			continue
		}

		endpointsSttingsParsed = append(endpointsSttingsParsed, endpointSetting)
	}
	endpointSettings := conf.EndpointSettings(endpointsSttingsParsed)

	if len(util.GetUnique(endpointSettings)) != len(endpointsSttingsParsed) {
		return nil, errors.New("duplicated endpoints found")
	}

	serviceSetting.Endpoints = endpointSettings

	serviceSetting.EndpointByID = map[string]*conf.EndpointSetting{}
	for _, endpointSetting := range serviceSetting.Endpoints {
		serviceSetting.EndpointByID[endpointSetting.GetID()] = endpointSetting
	}

	return serviceSetting, nil
}

func parseEndpointSetting(
	setting *conf.CoreSetting,
	serviceSetting *conf.ServiceSetting,
	itemConfig *config.Config,
) (*conf.EndpointSetting, error) {
	endpointSetting := &conf.EndpointSetting{}

	// open
	_, err := itemConfig.Get("open")
	if err != nil {
		endpointSetting.Open = true
	} else {
		opened, err := itemConfig.Bool("open")
		if err != nil {
			return nil, err
		}
		endpointSetting.Open = opened
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
	if len(util.Trim(rawExpose)) < 1 {
		return nil, errors.New("empty expose")
	}
	endpointSetting.Expose = util.Trim(rawExpose)
	endpointSetting.ExposeRegexp, err = regexp.Compile(endpointSetting.Expose)
	if err != nil {
		return nil, err
	}

	// source
	parsedSources, err := conf.ParseSource(itemConfig)
	if err != nil {
		return nil, err
	}
	sources := conf.MergeURLs(serviceSetting.Sources, parsedSources)
	if !conf.ValidateSources(sources) {
		return nil, errors.New("the sources in the defaults, must be valid url")
	}
	endpointSetting.Sources = sources

	// middleware
	middleware, err := parseMiddlewareSetting(setting, itemConfig, serviceSetting.Middleware)
	if err != nil {
		return nil, err
	}
	endpointSetting.Middleware = middleware

	return endpointSetting, nil
}

func parseEnvSetting(setting *conf.CoreSetting, itemConfig *config.Config) (*conf.EnvSetting, error) {
	_, err := itemConfig.Get("loglevel")
	if err != nil {
		return nil, err
	}
	logLevelInput, err := itemConfig.String("loglevel")

	var logLevel logging.Level
	if err == nil {
		logLevel, err = util.GetLogLevel(logLevelInput)
		if err != nil {
			return nil, err
		}
	}

	envSetting := conf.NewEnvSetting()
	envSetting.LogLevel = logLevel

	return envSetting, nil
}

func parseMiddlewareSetting(setting *conf.CoreSetting, itemConfig *config.Config, defaultName string) (conf.MiddlewareSetting, error) {
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

	if middleware, found := setting.Middlewares[middlewareName]; found {
		return middleware, nil
	}

	return nil, errors.New(fmt.Sprintf("invalid middleware found, `%s`", middlewareName))
}

func parseMiddlewares(itemConfig *config.Config) (conf.MiddlewaresSetting, error) {
	items, err := itemConfig.List("")
	if err != nil {
		return nil, err
	}

	middlewares := conf.MiddlewaresSetting{}

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

		middlewares[name] = conf.MiddlewareSetting{}
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
