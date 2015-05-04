package kasi

import "fmt"

func Run(yml string) {
	setting, err := ParseConfig(yml)
	if err != nil {
		panic(fmt.Sprintf("failed to read conf, %v", err))
	}

	HTTPServe(setting)
}
