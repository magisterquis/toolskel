// Program long_usage_line - Emit a long usage line
package main

/*
 * long_usage_line.go
 * Emit a long usage line
 * By J. Stuart McMurray
 * Created 20250326
 * Last Modified 20250326
 */

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() { os.Exit(rmain()) }
func rmain() int {
	/* Command-line flags. */
	var (
	/* TODO: Add flags. */
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %s [options]

Emit a long usage line.

In taberna quando sumus non curamus quid sit humus sed ad ludum properamus cui semper insudamus.

Options:
`,
			filepath.Base(os.Args[0]),
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* TODO: Meat and Potatoes. */

	return 0
}
