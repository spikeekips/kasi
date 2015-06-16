package markdown_config

import "encoding/json"

func ToJson(object interface{}) string {
	s, _ := json.MarshalIndent(object, "", "  ")
	return string(s)
}
