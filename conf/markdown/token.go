package markdown_config

import (
	"bufio"
	"errors"
	"fmt"
	"sort"
	"strings"
)

type TokenValue string

func (t TokenValue) String() string {
	return string(t)
}

func (t TokenValue) IsDelimiter() bool {
	return string(t) == string(ConfigValueDelimiter)
}

func (t TokenValue) IsValue() bool {
	return strings.HasPrefix(string(t), string(ConfigChar)) && strings.HasSuffix(string(t), string(ConfigChar))
}

func (t TokenValue) GetValue() (string, error) {
	if !t.IsValue() {
		return "", errors.New("is not config value")
	}

	return strings.TrimSuffix(strings.TrimPrefix(string(t), string(ConfigChar)), string(ConfigChar)), nil
}

type Token []TokenValue

func (s Token) String() string {
	return ToJson(s)
}

func (t *Token) Len() int {
	return len(*t)
}

func (t *Token) GetSource() TokenValue {
	return (*t)[0]
}

func (t *Token) NextTokenHasValue() bool {
	return (*t)[t.Len()-1].IsDelimiter()
}

func (t *Token) IndexDelimeter(start int) int {
	for n, i := range (*t)[start:] {
		if i.IsDelimiter() {
			return n
		}
	}

	return -1
}

func (t *Token) GetConfig(name string) Token {
	tokens := Token{}
	for _, i := range (*t)[1:] {
		if value, err := i.GetValue(); err == nil && value == name {
			tokens = append(tokens, i)
		}
	}

	return tokens
}

func (t *Token) GetConfigs() (config Config, err error) {
	// get delimiter
	var start int
	var indexDelimeters []int
	for {
		if start >= len(*t) {
			break
		}

		i := t.IndexDelimeter(start)
		if i < 1 {
			break
		}
		indexDelimeters = append(indexDelimeters, i)
		start += i + 1
	}

	config = Config{}

	if len(indexDelimeters) < 1 { // all TokenValue is the config without value
		for _, i := range *t {
			v, err := i.GetValue()
			if err != nil {
				continue
			}

			config[v] = []ConfigValue{}
		}
	}

	if len(config) < 1 {
		err = errors.New("no Config found")
	}

	return
}

func (t *Token) HasHeader() bool {
	return strings.Index(strings.TrimLeft(t.GetSource().String(), " \t"), ConfigHeader) >= 0
}

func (t *Token) IsRootHeader(name string) (found bool) {
	if !t.HasHeader() {
		return false
	}

	fmt.Println(t.GetConfig(name))
	config, err := t.GetConfigs()
	if err != nil {
		found = false
		return
	}

	_, found = config[name]
	fmt.Println(config, found)

	return
}

type Tokens map[int]Token

func (t *Tokens) GetKeys() []int {
	keys := []int{}
	for k, _ := range *t {
		keys = append(keys, k)
	}
	sorted := sort.IntSlice(keys)
	sort.Sort(sorted)

	return sorted
}

func Tokenize(s string) Tokens {
	var lineNo int

	tokens := Tokens{}

	addToToken := func(lineNo int, line string, token string) {
		if _, found := tokens[lineNo]; !found {
			tokens[lineNo] = Token{TokenValue(line)}
		}
		tokens[lineNo] = append(tokens[lineNo], TokenValue(token))
	}

	var inBlock bool

	scannerLine := bufio.NewScanner(strings.NewReader(s))
	for scannerLine.Scan() {
		line := scannerLine.Text()
		lineNo += 1

		// value list
		for _, b := range ConfigListChar {
			if strings.Trim(line, " \t") == string(b) {
				addToToken(lineNo, line, string(b))
				continue
			}
		}

		if indexBlock := strings.Index(line, ConfigChars); indexBlock >= 0 {
			blankLeft := len(strings.Trim(line[:indexBlock], " \t")) < 1
			blankRight := len(strings.Trim(line[indexBlock+len(ConfigChars):], " \t")) < 1
			if inBlock || (blankLeft && blankRight) {
				inBlock = !inBlock
				continue
			}
		}

		if inBlock {
			addToToken(lineNo, line, fmt.Sprintf("`%s`", line))
			continue
		}

		_, err := IndexByte([]byte(line), ConfigChar)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(strings.NewReader(line))
		scanner.Split(SplitToken)
		for scanner.Scan() {
			token := scanner.Text()
			if len(token) < 1 {
				continue
			}

			addToToken(lineNo, line, token)
		}
	}

	return tokens
}

func (t *Tokens) Objectify(rootName string, name string) {
	var configStarted bool
	for _, v := range *t {
		// is header
		if !configStarted && v.IsRootHeader(rootName) {
			fmt.Println(">>>>>>>>>>>", v[0])
			configStarted = true
		}
	}
}
