package util

import (
	"testing"

	"github.com/seanpont/assert"
)

func TestToJson(t *testing.T) {
	assert := assert.Assert(t)

	simpleStruct := struct {
		fieldHdden  string
		FieldString string
		FieldInt    int
		FieldArray  []string
	}{
		fieldHdden:  "fieldHdden",
		FieldString: "FieldString",
		FieldInt:    9,
		FieldArray: []string{
			"item0",
			"item1",
		},
	}

	cases := []struct {
		in   interface{}
		want string
	}{
		{
			simpleStruct, `{
  "FieldString": "FieldString",
  "FieldInt": 9,
  "FieldArray": [
    "item0",
    "item1"
  ]
}`,
		},
	}

	for _, c := range cases {
		assert.Equal(ToJson(c.in), c.want)
	}
}
