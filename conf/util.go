package kasi_conf

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/moraes/config"
	"github.com/spikeekips/kasi/util"
)

// Parse time unit, e.g. `3s`, `10m`.
// This follows the `time.ParseDuration()`, but also support day, month, year
// range such as `d`: day, `M`: Month, `y`: year .
func parseTimeUnit(s string) (time.Duration, error) {
	var err error
	var timeout_int64 int64
	var timeout_unit time.Duration = 1

	re_is_digit, _ := regexp.Compile(`^[\-\d][\d]*$`)
	matched := re_is_digit.FindStringSubmatch(s)
	if len(matched) == 1 {
		timeout_int64, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return time.Duration(0), errors.New(fmt.Sprintf("invalid time-unit format, `%v`", s))
		}
		return time.Duration(timeout_int64) * timeout_unit, nil
	}

	duration, err := time.ParseDuration(s)
	if err == nil {
		return duration, nil
	}

	re_time_unit, _ := regexp.Compile(`^(\d+)(s|m|h|d|M|y)$`)
	matched = re_time_unit.FindStringSubmatch(s)
	if len(matched) == 3 {
		timeout_int64, _ = strconv.ParseInt(matched[1], 10, 64)

		switch matched[2] {
		case "d":
			timeout_unit = time.Hour * time.Duration(24)
		case "M":
			timeout_unit = time.Hour * time.Duration(24*30)
		case "y":
			timeout_unit = time.Hour * time.Duration(24*365)
		}
	} else {
		return time.Duration(0), errors.New(fmt.Sprintf("invalid time-unit format, `%v`", s))
	}

	return time.Duration(timeout_int64) * timeout_unit, nil
}

func splitHostPort(s string) (*net.TCPAddr, error) {
	hostString, portString, err := net.SplitHostPort(s)
	if err != nil {
		return nil, err
	}

	if hostString == "0.0.0.0" {
		hostString = ""
	}

	ip := net.ParseIP(hostString)
	port, err := strconv.ParseInt(portString, 10, 64)
	if err != nil {
		return nil, err
	}

	return &net.TCPAddr{IP: ip, Port: int(port)}, nil
}

func joinPath(base string, target string) string {
	if path.IsAbs(target) {
		return target
	}
	return path.Join(base, target)
}

func joinURL(base string, target string) string {
	baseUrl, _ := url.Parse(base)
	targetUrl, _ := url.Parse(target)
	if len(targetUrl.Scheme) > 0 {
		return target
	}

	// join path
	baseUrl.Path = joinPath(baseUrl.Path, targetUrl.Path)

	// merge querystring
	queries := baseUrl.Query()
	for k, v := range targetUrl.Query() {
		queries[k] = v
	}
	baseUrl.RawQuery = queries.Encode()

	return UnescapeURL(baseUrl.String())
}

func isPath(maybePath string) bool {
	u, err := url.Parse(maybePath)
	if err != nil {
		return true
	}
	return len(u.Scheme) < 1
}

func UnescapeURL(s string) string {
	return strings.Replace(strings.Replace(s, "%7B", "{", -1), "%7D", "}", -1)
}

func MergeURLs(base []string, targets []string) (sources []string) {
	if len(targets) < 1 {
		return base
	}

	for _, t := range targets {
		if !isPath(t) || len(base) < 1 {
			sources = append(sources, t)
			continue
		}

		for _, b := range base {
			merged := joinURL(b, t)
			sources = append(sources, merged)
		}
	}

	return sources
}

func isValidURL(s string) bool {
	u, err := url.Parse(s)
	if err != nil {
		return false
	}

	if len(u.Scheme) < 1 {
		return false
	}

	return true
}

func validateSources(sources []string) bool {
	// sources must be valid url
	for _, i := range sources {
		if !isValidURL(i) {
			return false
		}
	}

	return true
}

func parseSource(itemConfig *config.Config) ([]string, error) {
	var err error

	_, err = itemConfig.String("source")
	if err == nil {
		return parseSourceString(itemConfig)
	}

	_, err = itemConfig.List("source")
	return parseSourceList(itemConfig)
}

func parseOneSource(s string) (string, error) {
	source := kasi_util.Trim(s)
	_, err := url.Parse(source)
	if err != nil {
		return "", err
	}

	return source, nil
}

func parseSourceString(itemConfig *config.Config) ([]string, error) {
	rawSource, err := itemConfig.String("source")
	if err != nil {
		return nil, err
	}

	var sources []string
	for _, i := range strings.Split(rawSource, ",") {
		if len(kasi_util.Trim(i)) < 1 {
			continue
		}
		sourceUrl, err := parseOneSource(i)
		if err != nil {
			return nil, err
		}
		sources = append(sources, sourceUrl)
	}

	return sources, nil
}

func parseSourceList(itemConfig *config.Config) ([]string, error) {
	rawSources, err := itemConfig.List("source")
	if err != nil {
		return nil, err
	}

	var sources []string
	for i := 0; i < len(rawSources); i++ {
		source, err := itemConfig.String(fmt.Sprintf("source.%d", i))
		if err != nil {
			return nil, err
		}

		if len(kasi_util.Trim(source)) < 1 {
			continue
		}
		sourceUrl, err := parseOneSource(source)
		if err != nil {
			return nil, err
		}
		sources = append(sources, sourceUrl)
	}

	return sources, nil
}
