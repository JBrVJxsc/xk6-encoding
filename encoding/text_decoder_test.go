package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//
// [WPT test]: https://github.com/web-platform-tests/wpt/blob/b5e12f331494f9533ef6211367dace2c88131fd7/encoding/textdecoder-labels.any.js
func TestTextDecoder(t *testing.T) {
	t.Parallel()

	ts := newTestSetup(t)
	err := executeTestScripts(ts, "./tests",
		"textdecoder-labels.js",
		"textdecoder-byte-order-marks.js",
	)
	assert.NoError(t, err)
}

func executeTestScripts(ts testSetup, base string, scripts ...string) error {
	for _, script := range scripts {
		program, err := compileFile(base, script)
		if err != nil {
			return err
		}

		if _, err = ts.rt.RunProgram(program); err != nil {
			return err
		}
	}

	return nil
}
