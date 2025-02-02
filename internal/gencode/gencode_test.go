package gencode

/*
 * gencode_test.go
 * Tests for gencode.go
 * By J. Stuart McMurray
 * Created 20250131
 * Last Modified 20250131
 */

import (
	"bytes"
	"embed"
	"encoding/json"
	"io/fs"
	"path"
	"reflect"
	"strings"
	"testing"
)

// testdata contains data used for testing.
//
//go:embed testdata
var testdata embed.FS

// testdataTopDir is the directory under which all other test data lives.
const testdataTopDir = "testdata"

// Make sure all of our templates load.  Since this happens on library load,
// we just need to make sure there's a test somewhere.
func TestInit(t *testing.T) {}

// Make sure we have a method for each template file.
func Test_HaveAllTmplFiles(t *testing.T) {
	/* Work out how many template functions we should have. */
	des, err := tmplFS.ReadDir(".")
	if nil != err {
		t.Fatalf("Error reading embedded template FS: %s", err)
	}
	want := len(des) + 1 /* For ToFile */

	/* Work out how many template functions we do have. */
	if got := reflect.TypeOf(Params{}).NumMethod(); got != want {
		t.Errorf(
			"Incorrect number of Template execution methods:\n"+
				" got: %d\n"+
				"want: %d\n"+
				"Re-run go generate?",
			got,
			want,
		)
	}
}

// getTestData gets the subFS of testdata corresponding to t.Name or calls
// t.Fatalf on error.
func getTestData(t *testing.T) fs.FS {
	/* Get this test's data subdirectory. */
	td, err := fs.Sub(testdata, path.Join(testdataTopDir, t.Name()))
	if nil != err {
		t.Fatalf(
			"Error geting subFS for %s: %s",
			t.Name(),
			err,
		)
	}
	/* Make sure we actually got something. */
	if _, err := td.Open("."); nil != err {
		t.Fatalf(
			"SubFS %s unusable: %s",
			t.Name(),
			err,
		)
	}
	return td
}

// Test file generation, both for correctness but also to make sure we have
// the same templates and methods.
func TestParams_Generation(t *testing.T) {
	var (
		haveFN = "have.json"
	)

	/* test tests that gen works as a generator using name as the name
	of the test cases.  It should be the Param.* method name. */
	testGen := func(t *testing.T, td fs.FS, name string, gen Generator) {
		/* Get the output we expect. */
		ns, err := fs.Glob(td, name+"_want.*")
		if nil != err {
			t.Fatalf("Error finding want file: %s", err)
		} else if 0 == len(ns) {
			t.Fatalf("No want file for %s", name)
		} else if 1 != len(ns) {
			t.Fatalf(
				"Multiple potential want files: %s",
				strings.Join(ns, " "),
			)
		}
		want, err := fs.ReadFile(td, ns[0])
		if nil != err {
			t.Fatalf("Error reading %s: %s", ns[0], err)
		}
		/* See if we get the right output. */
		if got, err := gen(); nil != err {
			t.Fatalf("Generation error: %s", err)
		} else if !bytes.Equal(got, want) {
			t.Errorf(
				"Incorrect generated file:\n"+
					" gotq: %q\n"+
					"wantq: %q\n"+
					"got:\n%s\n"+
					"want:\n%s\n",
				got,
				want,
				got,
				want,
			)
		}
	}

	/* testFromDir runs test in the subtest's subFS, which should be a
	directory containing a Params marshalled to json in test.json and
	file named X_want.EXT for each method X on params with an arbitrary
	EXT. */
	testFromDir := func(t *testing.T) {
		td := getTestData(t)
		/* Test params. */
		var p Params
		if b, err := fs.ReadFile(td, haveFN); nil != err {
			t.Fatalf("Error reading %s: %s", haveFN, err)
		} else if err := json.Unmarshal(b, &p); nil != err {
			t.Fatalf("Error unmarshalling Params: %s", err)
		}
		/* Try ALL the methods. */
		var (
			pv = reflect.ValueOf(p)
			pt = pv.Type()
		)
		for i := range pt.NumMethod() {
			name := pt.Method(i).Name
			/* If this isn't a generator, don't need to worry about
			it. */
			gen, ok := pv.Method(i).Interface().(Generator)
			if !ok {
				continue
			}
			t.Run(name, func(t *testing.T) {
				testGen(t, td, name, gen)
			})
		}
	}

	/* Test the test in each directory. */
	td := getTestData(t)
	des, err := fs.ReadDir(td, ".")
	if nil != err {
		t.Fatalf("Error getting test directories: %s", err)
	}
	for _, de := range des {
		t.Run(de.Name(), testFromDir)
	}
}
