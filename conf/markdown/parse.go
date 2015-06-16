package markdown_config

import (
	"bytes"
	"errors"
)

var ConfigHeader string = "#"
var ConfigChar byte = '`'
var ConfigChars string = "```"
var ConfigListChar []byte = []byte{'-', '*', '+'}
var ConfigValueDelimiter byte = ':'

func IndexByte(d []byte, s byte) (index int, err error) {
	index = bytes.IndexByte(d, s)
	if index < 0 {
		err = errors.New("not found")
		return
	}
	return
}

func SplitToken(data []byte, atEOF bool) (advance int, token []byte, err error) {
	start, err := IndexByte(data, ConfigChar)
	if err != nil {
		// name value delimeter
		for i := 0; i < len(data); i++ {
			if data[i] == ConfigValueDelimiter {
				return i + 1, data[i : i+1], nil
			}
		}

		return len(data), nil, nil
	}

	// name value delimeter
	for i := 0; i < start; i++ {
		if data[i] == ConfigValueDelimiter {
			return i + 1, data[i : i+1], nil
		}
	}

	// check multple chars
	lastIndex := start
	for i := start + 1; i < len(data); i++ {
		if data[i] == ConfigChar {
			if i == lastIndex+1 {
				lastIndex = i
			}
		}
	}

	if lastIndex != start { // got multple `ConfigChar`s
		// find closing
		cStart := lastIndex - start + 1
		lastIndexClose := lastIndex
		cClose := 0
		for i := lastIndex + 1; i < len(data); i++ {
			if data[i] == ConfigChar {
				if i == lastIndexClose+1 {
					cClose += 1
				} else {
					cClose = 1
				}
				lastIndexClose = i
			} else {
				cClose = 0
			}
			if cClose == cStart {
				break
			}
		}
		if lastIndex != lastIndexClose {
			return lastIndexClose + 1, nil, nil
		} else {
			return len(data), nil, nil
		}
	}

	// value list
	for i := 0; i < start; i++ {
		for _, b := range ConfigListChar {
			if _, err := IndexByte(data[i:i+1], b); err == nil {
				return i + 1, data[i : i+1], nil
			}
		}
	}

	for i := start + 1; i < len(data); i++ {
		if _, err := IndexByte(data[start+1:i+1], ConfigChar); err == nil {
			return i + 1, data[start : i+1], nil
		}
	}

	return len(data), nil, nil
}

func ParseConfig(s string) (Configs, error) {
	return Configs{}, nil
}
