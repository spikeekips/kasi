package kasi_util

import (
	"errors"
	"fmt"

	"github.com/op/go-logging"
)

var LogLevelNames = []string{
	"CRITICAL",
	"ERROR",
	"WARNING",
	"NOTICE",
	"INFO",
	"DEBUG",
}

func GetLogLevel(l string) (level logging.Level, err error) {
	for i, j := range LogLevelNames {
		if j == l {
			level = logging.Level(i)
			return
		}
	}

	err = errors.New(fmt.Sprintf("invalid loglevel name, `%s`", l))
	return
}
