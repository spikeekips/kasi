package kasi

import (
	"regexp"
	"strings"

	"github.com/moraes/config"
	"github.com/spikeekips/kasi/conf"
)

var RE_STRIP_TABED_LINE *regexp.Regexp = func() *regexp.Regexp { r, _ := regexp.Compile("^[\\s\\t]*$"); return r }()
var RE_STRIP_PREFIXED_TABS *regexp.Regexp = func() *regexp.Regexp { r, _ := regexp.Compile("^[\\t][\\t]*"); return r }()
var RE_STRIP_BLANKS *regexp.Regexp = func() *regexp.Regexp { r, _ := regexp.Compile("[\\s\\t]*$"); return r }()

// parse functions
func ParseConfig(yaml string) (*kasi_conf.CoreSetting, error) {
	// clean up the yaml string
	var cleanedYamlLine []string
	var cleanedYaml string
	for _, s := range strings.Split(yaml, "\n") {
		if RE_STRIP_TABED_LINE.Match([]byte(s)) {
			cleanedYamlLine = append(cleanedYamlLine, "")
			continue
		}

		striped := RE_STRIP_BLANKS.ReplaceAllString(s, "")
		cleanedYamlLine = append(cleanedYamlLine, striped)
	}
	cleanedYaml = strings.Join(cleanedYamlLine, "\n")

	conf, err := config.ParseYaml(cleanedYaml)
	if err != nil {
		return nil, err
	}

	setting, err := kasi_conf.ParseConfig(conf)
	if err != nil {
		return nil, err
	}

	return setting, nil
}
