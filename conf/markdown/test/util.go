package test_markdown_config

import (
	"io/ioutil"
	"os"
)

func loadFile(filename string) string {
	cur, _ := os.Getwd()
	f, _ := ioutil.ReadFile(cur + "/file/" + filename)
	return string(f)
}
