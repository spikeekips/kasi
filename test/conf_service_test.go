package kasi_t

import (
	"net"
	"testing"

	"github.com/seanpont/assert"
	"github.com/spikeekips/kasi/conf"
)

func TestGetServicesByBind(t *testing.T) {
	assert := assert.Assert(t)

	var coreSetting *kasi_conf.CoreSetting

	// blank services
	coreSetting = &kasi_conf.CoreSetting{}
	services := coreSetting.GetServicesByBind()
	assert.Equal(len(services), 0)

	makeCoreSettingsWithPorts := func(ports ...int) *kasi_conf.CoreSetting {
		services := []*kasi_conf.ServiceSetting{}
		for _, port := range ports {
			service := kasi_conf.ServiceSetting{
				Bind: &net.TCPAddr{Port: port},
			}
			services = append(services, &service)
		}

		return &kasi_conf.CoreSetting{
			Services: services,
		}
	}

	// 2 services, which have different bind
	coreSetting = makeCoreSettingsWithPorts(80, 90)

	services = coreSetting.GetServicesByBind()
	assert.Equal(len(services), 2)

	// if they have same bind
	coreSetting = makeCoreSettingsWithPorts(80, 80)

	services = coreSetting.GetServicesByBind()
	assert.Equal(len(services), 1)
}
