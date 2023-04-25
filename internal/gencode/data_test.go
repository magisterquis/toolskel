package gencode

/*
 * data_test.go
 * Tests for data.go
 * By J. Stuart McMurray
 * Created 20230418
 * Last Modified 20230425
 */

import (
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
