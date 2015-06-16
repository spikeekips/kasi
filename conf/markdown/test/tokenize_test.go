package test_markdown_config

import (
	"errors"
	"fmt"
	"sort"
	"testing"

	"github.com/seanpont/assert"
	"github.com/spikeekips/kasi/conf/markdown"
)

func check(input string, expect map[int][]string) (err error) {
	tokens := markdown_config.Tokenize(input)

	if len(tokens) != len(expect) {
		err = errors.New(fmt.Sprintf("mismatch, item count, `%v` != `%v`", len(tokens), len(expect)))
		return
	}

	for k, v := range expect {
		if _, found := tokens[k]; !found {
			err = errors.New(fmt.Sprintf("mismatch, missing line, `%v`", k))
			return
		}
		if len(v) != len(tokens[k])-1 {
			err = errors.New(fmt.Sprintf("mismatch, item count, `%v` != `%v`", len(tokens[k])-1, len(v)))
			return
		}
		for i, j := range v {
			if j != string(tokens[k][i+1]) {
				err = errors.New(fmt.Sprintf("mismatch, item[%d], `%v` != `%v`", i, tokens[k][i+1], j))
				return
			}
		}
	}

	return
}

func TestTokenize(t *testing.T) {
	assert := assert.Assert(t)

	tokens := markdown_config.Tokenize(loadFile("simple-default.md"))

	keys := []int{}
	for k, _ := range tokens {
		keys = append(keys, k)
	}
	sorted := sort.IntSlice(keys)
	sort.Sort(sorted)

	assert.Equal(len(tokens), 15)
	expected := []int{1, 6, 7, 8, 9, 10, 11, 12, 15, 19, 21, 23, 25, 27, 29}
	assert.Equal(len(sorted), len(expected))
	for n, i := range sorted {
		assert.Equal(i, expected[n])
	}
}

func TestTokenizeOneline(t *testing.T) {
	assert := assert.Assert(t)

	assert.Nil(check(
		"This is default `cache`: `yes`",
		map[int][]string{1: {"`cache`", ":", "`yes`"}},
	))
	assert.Nil(check(
		"This is default `cache`: show me `yes`",
		map[int][]string{1: {"`cache`", ":", "`yes`"}},
	))
	assert.Nil(check(
		"This is default cache`: `yes`",
		map[int][]string{1: {"`: `"}},
	))
	assert.Nil(check(
		"This is default c`ache`: `yes`",
		map[int][]string{1: {"`ache`", ":", "`yes`"}},
	))
	assert.Nil(check(
		"This is default `cache` `yes`",
		map[int][]string{1: {"`cache`", "`yes`"}},
	))
	assert.Nil(check(
		"This : is default `cache`: `yes`",
		map[int][]string{1: {":", "`cache`", ":", "`yes`"}},
	))
	assert.Nil(check(
		"This is `default` `cache`: `yes`",
		map[int][]string{1: {"`default`", "`cache`", ":", "`yes`"}},
	))
	assert.Nil(check(
		"This is default cache: yes",
		map[int][]string{},
	))
	assert.Nil(check(
		"This is default cache: `yes",
		map[int][]string{1: {":"}},
	))
	assert.Nil(check(
		"the entire configuration, the code expression, enclosed by `` ` `` will be interpreted for kasi.",
		map[int][]string{},
	))
	assert.Nil(check(
		"the entire `configuration`, the code expression, enclosed by `` ` `` will be interpreted for kasi.",
		map[int][]string{1: {"`configuration`"}},
	))
	assert.Nil(check(
		"the entire configuration, the code expression, enclosed by `` ` `` will be interpreted for `kasi`.",
		map[int][]string{1: {"`kasi`"}},
	))
	assert.Nil(check(
		"the entire configuration, ```the `code` expression```, enclosed by will be interpreted for kasi.",
		map[int][]string{},
	))
	assert.Nil(check(
		"configuration, ```the `code` expression``, enclosed`` by `will` : `be` interpreted for kasi.",
		map[int][]string{},
	))
	assert.Nil(check(
		"configuration, ```the `code` expression``, enclosed``` by `will` : `be` interpreted for kasi.",
		map[int][]string{1: {"`will`", ":", "`be`"}},
	))
	assert.Nil(check(
		"- This is `default` `cache`: `yes`",
		map[int][]string{1: {"-", "`default`", "`cache`", ":", "`yes`"}},
	))
	assert.Nil(check(
		"- This is `default-cache`: `yes`",
		map[int][]string{1: {"-", "`default-cache`", ":", "`yes`"}},
	))
}

func TestTokenizeMultiline(t *testing.T) {
	assert := assert.Assert(t)

	assert.Nil(check(
		"This is default `cache`: `yes`\n"+
			"This is default `cache`: `yes`\n"+
			"",
		map[int][]string{
			1: {"`cache`", ":", "`yes`"},
			2: {"`cache`", ":", "`yes`"},
		},
	))
	assert.Nil(check(
		"This is default `cache`: \n"+
			"```\n"+
			"findme\n"+
			"```",
		map[int][]string{
			1: {"`cache`", ":"},
			3: {"`findme`"},
		},
	))
	assert.Nil(check(
		"This is default `cache`: \n"+
			"\t``` \n"+
			"findme\n"+
			" ``` ",
		map[int][]string{
			1: {"`cache`", ":"},
			3: {"`findme`"},
		},
	))
	assert.Nil(check(
		"This is default `cache`: \n"+
			"\ta``` \n"+
			"findme\n"+
			" ```a",
		map[int][]string{
			1: {"`cache`", ":"},
		},
	))
	assert.Nil(check(
		"This is default `cache`: \n"+
			"\ta```showme \n"+
			"findme\n"+
			" ```a",
		map[int][]string{
			1: {"`cache`", ":"},
		},
	))
	assert.Nil(check(
		"This is default `cache`: ``` \n"+
			"findme\n"+
			" ```",
		map[int][]string{
			1: {"`cache`", ":"},
		},
	))
	assert.Nil(check(
		"This is default `cache`: ```: \n"+
			"findme\n"+
			" ```",
		map[int][]string{
			1: {"`cache`", ":"},
		},
	))
	assert.Nil(check(
		"This is default `cache`: `findme`\n"+
			" \n"+
			"This is default `list config`: \n"+
			"- `a`\n"+
			"- `b`",
		map[int][]string{
			1: {"`cache`", ":", "`findme`"},
			3: {"`list config`", ":"},
			4: {"-", "`a`"},
			5: {"-", "`b`"},
		},
	))
	assert.Nil(check(
		"This is default `cache`: \n"+
			"```\n"+
			"showme\n"+
			"showme\n"+
			"``` \n"+
			"This is default `list config`: \n"+
			"- `a`\n"+
			"- `b`",
		map[int][]string{
			1: {"`cache`", ":"},
			3: {"`showme`"},
			4: {"`showme`"},
			6: {"`list config`", ":"},
			7: {"-", "`a`"},
			8: {"-", "`b`"},
		},
	))
	assert.Nil(check(
		"This is default `cache`: \n"+
			"-\n"+
			"```\n"+
			"0showme\n"+
			"1showme\n"+
			"``` \n"+
			"-\n"+
			"```\n"+
			"0killme\n"+
			"1killme\n"+
			"``` \n",
		map[int][]string{
			1:  {"`cache`", ":"},
			2:  {"-"},
			4:  {"`0showme`"},
			5:  {"`1showme`"},
			7:  {"-"},
			9:  {"`0killme`"},
			10: {"`1killme`"},
		},
	))
}
