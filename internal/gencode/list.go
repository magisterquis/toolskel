package gencode

/*
 * list.go
 * List template types
 * By J. Stuart McMurray
 * Created 20230418
 * Last Modified 20230425
 */

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"

	"golang.org/x/exp/maps"
)

// descriptionTemplate is the subtemplate name which prints a description.
const descriptionTemplate = "description"

// ListTypes prints template types as a table to stdout.
func ListTypes() {
	/* Work out which templates we have available. */
	tns := maps.Keys(templates)
	sort.Strings(tns)

	/* Output will be nice and tabley. */
	tw := tabwriter.NewWriter(os.Stdout, 2, 8, 1, ' ', 0)
	defer tw.Flush()

	/* Get each template's description. */
	var buf bytes.Buffer
	for _, tn := range tns {
		buf.Reset()
		/* Get the subtemplate with the description. */
		t := templates[tn].Lookup(descriptionTemplate)
		if nil == t {
			log.Fatalf(
				"Template %q has no description subtemplate",
				tn,
			)
		}
		/* Extract the description from it. */
		if err := t.Execute(&buf, nil); nil != err {
			log.Fatalf(
				"Error getting description from "+
					"template %q: %s",
				tn,
				err,
			)
		}

		fmt.Fprintf(tw, "%s\t-\t%s\n", tn, buf.String())
	}
}
