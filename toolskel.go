// Program toolskel generates boilerplate for small tools written in Go.
package main

/*
 * toolskel.go
 * Generate command boilerplate
 * By J. Stuart McMurray
 * Created 20230204
 * Last Modified 20230427
 */

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/magisterquis/toolskel/internal/gencode"
)

func main() {
	var (
		noDate = flag.Bool(
			"no-date",
			false,
			"Do not set the Created/Modified date",
		)
		listTypes = flag.Bool(
			"list-types",
			false,
			"List available tool types",
		)
		tType = flag.String(
			"type",
			gencode.DefaultTType,
			"Tool `type` (see -list-types)",
		)
		author = flag.String(
			"author",
			defaultUsername(),
			"Author's `name`",
		)
		summaryCount = flag.Bool(
			"summary-count",
			false,
			"Generated code's summary prints a "+
				"completed task count ",
		)
		tagLog = flag.Bool(
			"tag-log",
			false,
			"Tag log output with argv[0]",
		)
		addVerbose = flag.Bool(
			"verbose-flag",
			false,
			"Add a -verbose flag",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %s [options] [toolname [tool description...]]

Generates boilerplate for a tool written in Go.

Options:
`,
			os.Args[0],
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* If we're just listing template types, life's easy. */
	if *listTypes {
		gencode.ListTypes()
		return
	}

	/* Fill in the rest of the data for the template. */
	data := gencode.Data{
		Name:         flag.Arg(0),
		Author:       *author,
		TagLog:       *tagLog,
		SummaryCount: *summaryCount,
		Verbose:      *addVerbose,
	}
	if "" != flag.Arg(1) {
		data.Description = strings.Join(flag.Args()[1:], " ")
	}
	if !*noDate {
		data.Today = time.Now().Format("20060102")
	}

	/* Generate the code itself. */
	if err := gencode.Generate(os.Stdout, *tType, data); nil != err {
		log.Fatalf("Error generating code: %s", err)
	}
}

// defaultUsername returns the current user's name or username, if available.
func defaultUsername() string {
	u, err := user.Current()
	if nil != err {
		log.Printf("Unable to get current user's info: %s", err)
		u = &user.User{}
	}
	for _, n := range []string{u.Name, u.Username} {
		if "" != n {
			return n
		}
	}
	return ""
}
