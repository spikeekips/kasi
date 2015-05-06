package conf

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"code.google.com/p/go-uuid/uuid"

	"github.com/spikeekips/kasi/util"
)

// EndpointSetting
type EndpointSetting struct {
	ID           string
	Open         bool
	Expose       string
	ExposeRegexp *regexp.Regexp
	Sources      []string
	Cache        time.Duration
	Timeout      time.Duration
	Middleware   MiddlewareSetting
}

func (setting *EndpointSetting) GetID() string {
	if setting.ID == "" {
		setting.ID = uuid.NewUUID().String()
	}

	return setting.ID
}

func (setting *EndpointSetting) String() string {
	return util.ToJson(setting)
}

type EndpointSettings []*EndpointSetting

func (setting EndpointSettings) Len() int {
	return len(setting)
}

func (setting EndpointSettings) AreEqual(a int, b int) bool {
	if setting[a] == nil || setting[b] == nil {
		return false
	}

	return setting[a].Expose == setting[b].Expose
}

func (setting *EndpointSetting) GetTargetURL(url url.URL) (string, error) {
	subNames := setting.ExposeRegexp.SubexpNames()
	subMatches := setting.ExposeRegexp.FindAllStringSubmatch(url.Path, -1)
	if len(subMatches) < 1 {
		return "", errors.New("failed to get pattern")
	}

	targetPattern := string([]byte(setting.Sources[0]))
	for i, n := range subMatches[0] {
		targetPattern = strings.Replace(
			strings.Replace(
				targetPattern,
				fmt.Sprintf("{%s}", subNames[i]),
				n,
				-1,
			),
			fmt.Sprintf("{%d}", i-1),
			n,
			-1,
		)
	}

	return targetPattern, nil
}
