package test_yaml_conf

import (
	"testing"

	"github.com/seanpont/assert"
	"github.com/spikeekips/kasi"
)

func TestMiddleware(t *testing.T) {
	assert := assert.Assert(t)

	yml := loadFile("conf_middlewares.yml")
	setting, err := kasi.ParseConfig(yml)
	assert.Equal(err, nil)

	assert.Equal(len(setting.Middlewares), 2)
	service0, _ := setting.GetServiceByID("service0")
	assert.Equal(service0.Middleware, "base")
	service1, _ := setting.GetServiceByID("service1")
	assert.Equal(service1.Middleware, "extend")

	assert.Equal(len(service0.Endpoints[0].Middleware), 2)
	assert.Equal(service0.Endpoints[0].Middleware[0], "a.js")
	assert.Equal(service0.Endpoints[0].Middleware[1], "b.js")

	assert.Equal(len(service0.Endpoints[0].Middleware), 2)
	assert.Nil(service0.Endpoints[1].Middleware)

	assert.Equal(len(service1.Endpoints[0].Middleware), 2)
	assert.Equal(service1.Endpoints[0].Middleware[0], "a.js")
	assert.Equal(service1.Endpoints[0].Middleware[1], "b.js")

	assert.Equal(len(service1.Endpoints[1].Middleware), 2)
	assert.Equal(service1.Endpoints[1].Middleware[0], "c.js")
	assert.Equal(service1.Endpoints[1].Middleware[1], "d.js")
}
