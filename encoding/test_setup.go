package encoding

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/grafana/sobek"
	"github.com/stretchr/testify/require"
	"go.k6.io/k6/js/modulestest"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/metrics"
	"gopkg.in/guregu/null.v3"
)

// testSetup is a helper struct holding components
// necessary to test the redis client, in the context
// of the execution of a k6 script.
type testSetup struct {
	rt      *sobek.Runtime
	state   *lib.State
	samples chan metrics.SampleContainer
}

// newTestSetup initializes a new test setup.
// It prepares a test setup with a mocked redis server and a sobek runtime,
// ready to execute scripts as if being executed in the
// main context of k6.
func newTestSetup(t testing.TB) testSetup {
	rt := sobek.New()
	rt.SetFieldNameMapper(sobek.TagFieldNameMapper("json", true))

	// We compile the Web Platform testharness script into a sobek.Program
	encodingsProgram, err := compileFile("./tests/resources", "encodings.js")
	require.NoError(t, err)

	// We execute the harness script in the sobek runtime
	// in order to make the Web Platform assertion functions available
	// to the tests.
	_, err = rt.RunProgram(encodingsProgram)
	require.NoError(t, err)

	// We compile the Web Platform testharness script into a sobek.Program
	assertProgram, err := compileFile("./tests/utils", "assert.js")
	require.NoError(t, err)

	// We execute the harness script in the sobek runtime
	// in order to make the Web Platform assertion functions available
	// to the tests.
	_, err = rt.RunProgram(assertProgram)
	require.NoError(t, err)

	samples := make(chan metrics.SampleContainer, 1000)

	state := &lib.State{
		Options: lib.Options{
			SystemTags: metrics.NewSystemTagSet(
				metrics.TagURL,
				metrics.TagProto,
				metrics.TagStatus,
				metrics.TagSubproto,
			),
			UserAgent: null.StringFrom("TestUserAgent"),
		},
		Samples:        samples,
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(metrics.NewRegistry()),
		Tags:           lib.NewVUStateTags(metrics.NewRegistry().RootTagSet()),
	}

	vu := &modulestest.VU{
		RuntimeField: rt,
		StateField:   state,
	}

	m := new(RootModule).NewModuleInstance(vu)
	require.NoError(t, rt.Set("TextDecoder", m.Exports().Named["TextDecoder"]))

	return testSetup{
		rt:      rt,
		state:   state,
		samples: samples,
	}
}

// compileFile compiles a javascript file as a sobek.Program.
func compileFile(base, name string) (*sobek.Program, error) {
	fname := path.Join(base, name)

	//nolint:forbidigo
	f, err := os.Open(filepath.Clean(fname))
	if err != nil {
		return nil, err
	}
	defer func() {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	str := string(b)
	program, err := sobek.Compile(name, str, false)
	if err != nil {
		return nil, err
	}

	return program, nil
}
