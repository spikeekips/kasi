package kasi_conf

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"time"

	"code.google.com/p/go-uuid/uuid"

	"github.com/spikeekips/kasi/util"
)

// ServiceSetting
type ServiceSetting struct {
	ID           string
	Hostnames    []string
	Bind         *net.TCPAddr
	SSL          *SSLSetting
	Endpoints    EndpointSettings
	EndpointByID map[string]*EndpointSetting

	// default values for each endpoints
	Sources []string

	Middleware string

	Timeout time.Duration
	Cache   time.Duration
}

func (setting *ServiceSetting) GetID() string {
	if setting.ID == "" {
		setting.ID = uuid.NewUUID().String()
	}

	return setting.ID
}

func (setting *ServiceSetting) String() string {
	return kasi_util.ToJson(setting)
}

func (setting *ServiceSetting) Opened() bool {
	if len(setting.Endpoints) < 1 {
		return false
	}
	for _, endpointSetting := range setting.Endpoints {
		if endpointSetting.Opened() {
			return true
		}
	}

	return false
}

func (setting *ServiceSetting) GetPatterns() []string {
	portExpression := ""
	// the builtin `http` does not understand the default port.
	if sort.SearchInts([]int{80, 443}, setting.Bind.Port) == 2 {
		portExpression = fmt.Sprintf(":%d", setting.Bind.Port)
	}

	var exposes []string
	for _, s := range setting.Hostnames {
		expose := fmt.Sprintf(
			"%s%s/",
			kasi_util.RStripSlash(s),
			portExpression,
		)
		exposes = append(exposes, expose)
	}

	return exposes
}

func (setting *ServiceSetting) GetMatchedEndpoint(path string) (*EndpointSetting, error) {
	for _, endpoint := range setting.Endpoints {
		matched := endpoint.exposeRegexp.Match([]byte(path))
		if !matched {
			continue
		}
		return endpoint, nil
	}
	return nil, errors.New("failed to find endpoint")
}
