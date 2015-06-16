package test_markdown_config

import (
	"testing"

	"github.com/spikeekips/kasi/conf/markdown"
)

func TestObjectify(t *testing.T) {
	//assert := assert.Assert(t)

	var doc string

	doc = "# this is `kasi config`\n" +
		"This is default `cache`: `yes`\n" +
		""
	tokens := markdown_config.Tokenize(doc)
	tokens.Objectify("kasi config", "")
}
