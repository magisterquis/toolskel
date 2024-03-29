// Program cooltool - A cool program
package main

/*
 * cooltool.go
 * A cool program
 * By MysteryDev
 * Created in the past
 * Last Modified in the past
 */

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	/* ProgramStart notes when the program has started for printing the
	elapsed time when the program ends. */
	ProgramStart = time.Now()
)

func main() {
	/* Tag log messages with argv[0]. */
	log.SetPrefix("[" + os.Args[0] + "] ")

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

A cool program

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
			"Done in %s.",
			time.Since(ProgramStart).Round(time.Millisecond),
		)
	}
}
