package kasi_t

import (
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
)

func loadFile(filename string) string {
	cur, _ := os.Getwd()
	f, _ := ioutil.ReadFile(cur + "/file/" + filename)
	return string(f)
}

func GetUrl(s string) *url.URL {
	url, _ := url.Parse(s)
	return url
}

func GetRegexp(s string) *regexp.Regexp {
	re, _ := regexp.Compile(s)
	return re
}
