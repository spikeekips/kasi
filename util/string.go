package kasi_util

import (
	"encoding/json"
	"regexp"
	"sort"
	"strings"
)

func Trim(s string) string {
	return strings.Trim(s, " \t")
}

var RE_RSTRIP_SLASH *regexp.Regexp

func RStripSlash(h string) string {
	if RE_RSTRIP_SLASH == nil {
		RE_RSTRIP_SLASH, _ = regexp.Compile(`[\/]*$`)
	}

	return RE_RSTRIP_SLASH.ReplaceAllString(h, "")
}

func ToJson(object interface{}) string {
	s, _ := json.MarshalIndent(object, "", "  ")
	return string(s)
}

func InArray(s []string, x string) bool {
	l := sort.StringSlice(s)
	sort.Sort(l)

	idx := sort.SearchStrings(l, x)
	if idx == len(s) {
		return false
	}

	return l[idx] == x
}
