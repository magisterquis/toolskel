package main

/*
 * gengencode_test.go
 * Make sure this thing works
 * By J. Stuart McMurray
 * Created 20250130
 * Last Modified 20250130
 */

import (
	"bytes"
	"embed"
	"io/fs"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

//go:embed testdata
var testdata embed.FS

func TestDo(t *testing.T) {
	const (
		wantGoFn  = "want.go"
		wantErrFn = "wanterr"
	)

	// Deterministic date.
	date := time.Date(2525, time.December, 25, 13, 37, 00, 00, time.UTC)
	// Try ALL the test directories!
	dn := "testdata/do"
	des, err := fs.ReadDir(testdata, dn)
	if nil != err {
		t.Fatalf("Error reading testdata directory %s: %s", dn, err)
	}
	for _, de := range des {
		t.Run(de.Name(), func(t *testing.T) {
			/* Put our test in its own directory. */
			td := t.TempDir()
			sd := path.Join(dn, de.Name())
			sfs, err := fs.Sub(testdata, sd)
			if nil != err {
				t.Fatalf(
					"Error getting testdata subds %s: %s",
					sd,
					err,
				)
			}
			if err := os.CopyFS(td, sfs); nil != err {
				t.Fatalf("Error unpacking test data: %s", err)
			}
			/* Work out what we'd like. */
			want, err := fs.ReadFile(sfs, wantGoFn)
			if nil != err {
				t.Fatalf("Error reading %s: %s", wantGoFn, err)
			}
			b, err := fs.ReadFile(sfs, wantErrFn)
			if nil != err {
				t.Fatalf("Error reading %s: %s", wantErrFn, err)
			}
			wantErr := strings.TrimSpace(string(b))
			/* Generate an output file. */
			got, err := do(td, date)
			/* See if it all worked. */
			if !bytes.Equal(got, want) {
				t.Errorf(
					"Incorrect generated code:\n"+
						"got:\n%s\n"+
						"want:\n%s",
					got,
					want,
				)
			}
			var gotErr string
			if nil != err {
				gotErr = err.Error()
			}
			if gotErr != string(wantErr) {
				t.Errorf(
					"Incorrect error:\n"+
						" got: %s\n"+
						"want: %s",
					gotErr,
					wantErr,
				)
			}
		})
	}
}
