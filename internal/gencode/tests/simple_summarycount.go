// Program cooltool: A cool program
package main

/*
 * cooltool.go
 * cooltool: A cool program
 * By MysteryDev
 * Created in the past
 * Last Modified in the past
 */

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"
)

var (
	/* ProgramStart notes when the program has started for printing the
	elapsed time when the program ends. */
	ProgramStart = time.Now()

	/* NDone keeps track of the number of things we've done. */
	NDone atomic.Uint64
)

func main() {
	/* Command-line flags. */
	var (
		noSummary = flag.Bool(
			"no-summary",
			false,
			"Don't print a summary on exit",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %s [options]

cooltool: A cool program

Options:
`,
			os.Args[0],
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* TODO: Meat and Potatoes. */

	/* All done. */
	if !*noSummary {
		log.Printf(
			"Done.  Finished %d in %s.",
			NDone.Load(),
			time.Since(ProgramStart).Round(time.Millisecond),
		)
	}
}
