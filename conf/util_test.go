package kasi_conf

import (
	"net"
	"testing"
	"time"

	"github.com/seanpont/assert"
)

func TestParseTimeUnit(t *testing.T) {
	assert := assert.Assert(t)

	cases := []struct {
		in   string
		want time.Duration
	}{
		{"10", time.Nanosecond * 10},
		{"100000", time.Nanosecond * 100 * 1000},
		{"1s", time.Second * 1},
		{"2s", time.Second * 2},
		{"3m", time.Minute * 3},
		{"4h", time.Hour * 4},
		{"5d", time.Hour * 24 * 5},
		{"6M", time.Hour * 24 * 30 * 6},
		{"7y", time.Hour * 24 * 365 * 7},
	}
	for _, c := range cases {
		assert.Equal(
			func() time.Duration {
				p, err := parseTimeUnit(c.in)
				assert.Nil(err)
				return p
			}(),
			c.want,
		)
	}
}

func TestSplitHostPort(t *testing.T) {
	assert := assert.Assert(t)

	cases := []struct {
		in   string
		want *net.TCPAddr
	}{
		{":80", &net.TCPAddr{Port: 80}},
		{"0.0.0.0:80", &net.TCPAddr{Port: 80}},
		{"*:80", &net.TCPAddr{Port: 80}},
		{"192.168.0.1:80", &net.TCPAddr{IP: net.ParseIP("192.168.0.1"), Port: 80}},
	}
	for _, c := range cases {
		assert.Equal(
			func() *net.TCPAddr {
				p, err := splitHostPort(c.in)
				assert.Nil(err)
				return p
			}(),
			c.want,
		)
	}
}

func TestJoinURL(t *testing.T) {
	assert := assert.Assert(t)

	var base, target string

	// 2 FQDN url
	base = "http://a0.com/b0/c0.html"
	target = "http://a1.com/b1/c1.html"
	assert.Equal(joinURL(base, target), target)

	// 1 FQDN and 1 path
	base = "http://a0.com/b0/c0.html"
	target = "/b1/c1.html"
	assert.Equal(joinURL(base, target), "http://a0.com"+target)

	// 1 invalid url and 1 valid url
	base = "****************\\*"
	target = "/b1/c1.html"
	assert.Equal(joinURL(base, target), target)

	// 2 path, but one is relative
	base = "/b0/c0.html"
	target = "b1/c1.html"
	assert.Equal(joinURL(base, target), base+"/"+target)

	// merge querystring
	base = "/b0/c0.html?a=1&b=2"
	target = "b1/c1.html?a=3"
	assert.Equal(joinURL(base, target), "/b0/c0.html/b1/c1.html?a=3&b=2")
}

func TestIsPath(t *testing.T) {
	assert := assert.Assert(t)

	cases := []struct {
		in   string
		want bool
	}{
		{"/a/", true},
		{"://a/", true},
		{"http://a.com/", false},
		{"p://a.com/", false},
	}
	for _, c := range cases {
		assert.Equal(
			isPath(c.in),
			c.want,
		)
	}
}
