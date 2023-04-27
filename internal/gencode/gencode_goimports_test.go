//go:build testgoimports

package gencode

/*
 * gencode_goimports_test.go
 * GoImports Tests for gencode.go
 * By J. Stuart McMurray
 * Created 20230415
 * Last Modified 20230427
 */

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestMainGoImports(t *testing.T) {
	/* Make sure we actually have goimports. */
	gidir := t.TempDir()
	cmd := exec.Command(
		"go",
		"install",
		"golang.org/x/tools/cmd/goimports@latest",
	)
	cmd.Env = append(cmd.Environ(), fmt.Sprintf("GOBIN=%s", gidir))
	o, err := cmd.CombinedOutput()
	if nil != err {
		if 0 == len(o) {
			t.Fatalf("Error fetching goimports: %s", err)
		}
		t.Fatalf("Error fetching goimports: %s\n%s", err, o)
	}
	goimports := filepath.Join(gidir, "goimports")

	/* function to call the above-downloaded goimports and format a source
	file. */
	runGoImports := func(sCode []byte) ([]byte, error) {
		cmd := exec.Command(goimports)
		cmd.Stdin = bytes.NewReader(sCode)
		return cmd.CombinedOutput()
	}

	/* Test ALL the things. */
	for _, c := range TestCases {
		c := c /* Until someone fixes it :| */
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			/* Generate the code with this config. */
			var buf bytes.Buffer
			if err := Generate(&buf, c.tType, c.data); nil != err {
				t.Errorf("generating code: %s", err)
				return
			}
			sCode := buf.Bytes()
			/* Run it through goimports. */
			iCode, err := runGoImports(sCode)
			if nil != err {
				if 0 == len(iCode) {
					t.Errorf(
						"goimports returned error: %s",
						err,
					)
					return
				}
				t.Errorf(
					"goimports returned error: %s\n%s",
					err,
					iCode,
				)
				return
			}
			/* Tell someone if it's not the same. */
			errorIfDiff(t, sCode, iCode, "generated", "goimported")
		})
	}
}
