package gencode

/*
 * data_test.go
 * Tests for data.go
 * By J. Stuart McMurray
 * Created 20230418
 * Last Modified 20240419
 */

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestDataImportsBlock(t *testing.T) {
	for _, c := range []struct {
		Have []string
		Want string
	}{{
		Want: `import ()`,
	}, {
		Have: []string{"flag", "os", "fmt"},
		Want: `import (
	"flag"
	"fmt"
	"os"
)`,
	}, {
		Have: []string{
			"golang.org/x/exp/maps",
			"net",
			"github.com/test/001",
			"math",
			"github.com/example/002",
		},
		Want: `import (
	"math"
	"net"

	"github.com/example/002"
	"github.com/test/001"
	"golang.org/x/exp/maps"
)`,
	}} {
		c := c /* :S */
		t.Run(strings.Join(c.Have, " "), func(t *testing.T) {
			t.Parallel()
			var d Data
			if nil == c.Have {
				c.Have = make([]string, 0)
			}
			sort.Strings(c.Have)
			d.SetDefaults()
			got := d.WithImports(c.Have...).ImportsBlock()
			if got != c.Want {
				t.Errorf("got:\n%s", got)
			}
		})
	}
}

func TestDataSetDefaults(t *testing.T) {
	var data Data
	data.SetDefaults()

	wd, err := os.Getwd()
	if nil != err {
		t.Fatalf("Error getting current directory: %s", err)
	}
	want := filepath.Base(wd)

	if got := data.Name; got != want {
		t.Errorf(
			"Incorrect default name\n"+
				" cwd: %s\n"+
				" got: %s\n"+
				"want: %s\n",
			wd,
			got,
			want,
		)
	}
}
