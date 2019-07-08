package jq

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type jqcase struct {
	filter string
	input  string
	output string
	err    string
}

var cases = []jqcase{
	jqcase{
		filter: ".",
		input:  `{"some":"json"}`,
		output: string("{\n  " + `"some": "json"` + "\n}\n"),
		err:    "",
	},
	jqcase{
		filter: ".some",
		input:  `{"some":"json"}`,
		output: string(`"json"` + "\n"),
		err:    "",
	},
	jqcase{
		filter: "abc",
		input:  `{"some":"json"}`,
		output: "",
		err:    "jq: error: abc/0 is not defined at <top-level>, line 1:\nabc\njq: 1 compile error\n",
	},
	jqcase{
		filter: ".",
		input:  `"some":"json"}`,
		output: string(`"some"` + "\n"),
		err:    "parse error: Expected string key before ':' at line 1, column 7\n",
	},
}

func TestEvalStream(t *testing.T) {
	for _, e := range cases {
		var jqr *JQ
		Init()
		Path = "/usr/local/bin/jq"
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		jqr = &JQ{
			J: "",
			Q: e.filter,
		}
		var obuff bytes.Buffer
		var ebuff bytes.Buffer
		input := strings.NewReader(e.input)

		if err := jqr.EvalStream(ctx, input, &obuff, &ebuff); err != nil {
			fmt.Println(err.Error())
		}

		assert.Equal(t, e.output, obuff.String())
		assert.Equal(t, e.err, ebuff.String())
	}
}

func TestEval(t *testing.T) {
	for _, e := range cases {
		var jqr *JQ
		Init()
		Path = "/usr/local/bin/jq"
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		jqr = &JQ{
			J: e.input,
			Q: e.filter,
		}
		var obuff bytes.Buffer
		var ebuff bytes.Buffer

		if err := jqr.Eval(ctx, &obuff, &ebuff); err != nil {
			fmt.Println(err.Error())
		}

		assert.Equal(t, e.output, obuff.String())
		assert.Equal(t, e.err, ebuff.String())
	}
}

func TestOptions(t *testing.T) {
		var jqr *JQ
		Init()
		Path = "/usr/local/bin/jq"
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		jqr = &JQ{
			J: `{"some":"json"}`,
			Q: ".",
			O: []JQOpt{
				JQOpt{
					Name: "monochrome-output",
					Enabled: true,
				},
				JQOpt{
					Name: "compact-output",
					Enabled: true,
				},
			},
		}
		var obuff bytes.Buffer
		var ebuff bytes.Buffer

		if err := jqr.Eval(ctx, &obuff, &ebuff); err != nil {
			fmt.Println(err.Error())
		}

		assert.Equal(t, "{\"some\":\"json\"}\n", obuff.String())
		assert.Equal(t, "", ebuff.String())
}
