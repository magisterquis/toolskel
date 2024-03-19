package gencode

/*
 * gencode_test.go
 * Tests for gencode.go
 * By J. Stuart McMurray
 * Created 20230415
 * Last Modified 20240419
 */

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// testWants contains the contents of the tests/ directory
//
//go:embed tests/*.go
var testWants embed.FS

// testWantsDir is the directory in testWants with the files for TestCases.
const testWantsDir = "tests"

// TestCases are common test cases for various tests.  Each test should have
// a corresponding file in tests/ named after .name, with /'s replaced by
// _'s and a .go suffix.
var TestCases = []struct {
	name  string
	tType string /* Template name, less template/ and .tmpl */
	data  Data
	want  []byte
}{{
	name: "simple.go",
}, {
	name: "simple/summarycount.go",
	data: Data{
		SummaryCount: true,
	},
}, {
	name: "simple/taglog.go",
	data: Data{
		TagLog: true,
	},
}, {
	name: "simple/verbose.go",
	data: Data{
		Verbose: true,
	},
}, {
	name: "simple/summarycountverbose.go",
	data: Data{
		SummaryCount: true,
		Verbose:      true,
	},
}, {
	name:  "parallel.go",
	tType: "parallel",
}, {
	name:  "parallel/summarycount.go",
	tType: "parallel",
	data: Data{
		SummaryCount: true,
	},
}, {
	name:  "parallel/verbose.go",
	tType: "parallel",
	data: Data{
		Verbose: true,
	},
}, {
	name: "library.go",
	data: Data{
		Name: "main", /* For testing. */
	},
	tType: "library",
}}

// init populates TestCases's data fields.
func init() {
	var err error
	for i, c := range TestCases {
		/* Gotta have a name. */
		if "" == c.name {
			panic(fmt.Sprintf("no name for test %d", i))
		}

		/* Get the file. */
		fn := filepath.Join(
			testWantsDir,
			strings.Replace(c.name, "/", "_", -1),
		)
		c.want, err = testWants.ReadFile(fn)
		if nil != err {
			panic(fmt.Sprintf("reading test file %s: %s", fn, err))
		}
		/* Save it back. */
		TestCases[i] = c
	}
}

func TestGenCode(t *testing.T) {
	for _, c := range TestCases {
		c := c /* :( */
		/* If there's no known good for this one, skip it. */
		if 0 == len(c.want) {
			continue
		}
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			if err := Generate(&buf, c.tType, c.data); nil != err {
				t.Errorf("error: %s", err)
				return
			}
			errorIfDiff(t, buf.Bytes(), c.want, "", "")
		})
	}
}

// TestWantBuild tests that the test cases' known-goods actually build.
func TestWantBuild(t *testing.T) {
	des, err := testWants.ReadDir(testWantsDir)
	if nil != err {
		t.Fatalf("Error reading embedded FS: %s", err)
	}

	for _, de := range des {
		de := de /* D: */
		t.Run(de.Name(), func(t *testing.T) {
			t.Parallel()
			if de.IsDir() {
				t.Errorf("Got a directory, expected a file")
				return
			}
			efn := filepath.Join(testWantsDir, de.Name())
			b, err := testWants.ReadFile(efn)
			if nil != err {
				t.Errorf("Error reading file %s: %s", efn, err)
				return
			}
			td := t.TempDir()
			fn := filepath.Join(td, de.Name())
			if err := os.WriteFile(fn, b, 0660); nil != err {
				t.Errorf("Error writing to %s: %s", fn, err)
				return
			}
			if _, err := combinedOutput(
				t,
				td,
				"go mod init tstest",
			); nil != err {
				t.Errorf("Error adding go.mod: %s", err)
				return
			}
			if _, err := combinedOutput(
				t,
				td,
				"go run . -h",
			); nil != err {
				t.Errorf("Build failed with error: %s", err)
				return
			}
		})
	}
}

// combinedOutputError is returned by combinedOutput when the underlying
// exec.Cmd.CombinedOutput returns an error.
type combinedOutputError struct {
	Output []byte
	Err    error
}

// Unwrap returns the underlying error but no output.
func (err combinedOutputError) Unwrap() error { return err.Err }

// Error implements the error interface.
func (err combinedOutputError) Error() string {
	if 0 == len(err.Output) {
		return err.Error()
	}
	return fmt.Sprintf("%s\nOutput:\n%s)", err.Err, err.Output)
}

// combinedOutput returns the output of running the command in a shell in the
// given directory, plus any errors encountered.  If exec.Cmd.CombinedOutput
// returns an error, combinedOutput returns a combinedOutputError.
func combinedOutput(t *testing.T, dir, cmd string) ([]byte, error) {
	t.Helper()
	if 0 == len(cmd) {
		return nil, fmt.Errorf("empty command")
	}
	/* Yeah, Windows.  PRs welcome. */
	c := exec.Command("/bin/sh", "-c", cmd)
	c.Dir = dir
	o, err := c.CombinedOutput()
	if nil != err {
		return nil, combinedOutputError{Output: o, Err: err}
	}
	return o, nil
}
