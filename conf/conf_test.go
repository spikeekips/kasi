package kasi_conf

/*

import (
	"net"
	"strconv"
	"testing"

	"github.com/seanpont/assert"
	"github.com/spikeekips/kasi/test"
)

func TestGetExposeForBind(t *testing.T) {
	assert := assert.Assert(t)

	var serviceSetting *ServiceSetting
	var hostname string
	var port int

	// with official port, 80
	hostname = "a0.dev"
	port = 80
	serviceSetting = &ServiceSetting{
		Hostnames: []string{hostname},
		Bind:      &net.TCPAddr{Port: port},
		Endpoints: EndpointSettings{
			&EndpointSetting{
				exposeRegexp: kasi_t.GetRegexp("^/(?P<filename>.*)/$"),
				SourceUrl:    kasi_t.GetUrl(RStripSlash("http://localhost:8080/")),
				Source:       "/f/{filename}",
			},
		},
	}

	// with unofficial port
	hostname = "a0.dev"
	port = 8080
	serviceSetting = &ServiceSetting{
		Hostnames: []string{hostname},
		Bind:      &net.TCPAddr{Port: port},
		Endpoints: EndpointSettings{
			&EndpointSetting{
				exposeRegexp: kasi_t.GetRegexp("^/(?P<filename>.*)/$"),
				SourceUrl:    kasi_t.GetUrl(RStripSlash("http://localhost:8080/")),
				Source:       "/f/{filename}",
			},
		},
	}
	assert.Equal(serviceSetting.GetExposeForBind()[0], hostname+":"+strconv.Itoa(port)+"/")
}

func TestGetMatchedEndpoint(t *testing.T) {
	assert := assert.Assert(t)

	endpointSetting := &EndpointSetting{
		exposeRegexp: kasi_t.GetRegexp("^/prefix/(?P<filename>.*)/$"),
		SourceUrl:    kasi_t.GetUrl(RStripSlash("http://localhost:8080/")),
		Source:       "/f/{filename}",
	}

	serviceSetting := ServiceSetting{
		Hostnames: []string{"a.dev"},
		Bind:      &net.TCPAddr{Port: 8080},
		Endpoints: EndpointSettings{
			endpointSetting,
		},
	}
	_, err := serviceSetting.GetMatchedEndpoint("/")
	assert.NotNil(err)

	_, err = serviceSetting.GetMatchedEndpoint("/prefix/")
	assert.NotNil(err)

	_, err = serviceSetting.GetMatchedEndpoint("/prefix/will-not-match-without-trailing-slash")
	assert.NotNil(err)

	matchedEndpoint, err := serviceSetting.GetMatchedEndpoint("/prefix/will-be-matched/")
	assert.Nil(err)
	assert.Equal(matchedEndpoint, endpointSetting)
	assert.Equal(ToJson(matchedEndpoint), ToJson(endpointSetting))
}
*/
