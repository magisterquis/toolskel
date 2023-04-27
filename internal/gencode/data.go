package gencode

/*
 * data.go
 * Data we pass to templates
 * By J. Stuart McMurray
 * Created 20230418
 * Last Modified 20230427
 */

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"golang.org/x/exp/maps"
)

// The following default values are compile-time settable.
var (
	defaultDescription = "A cool program"
	defaultProgramName = "cooltool"
	defaultAuthorName  = "MysteryDev"
)

// Data is used to pass data to the template being executed.
type Data struct {
	/* Strings which go right into code. */
	Name         string              /* Package name. */
	Description  string              /* Short description. */
	Author       string              /* Author's name. */
	Today        string              /* Curent date. */
	SummaryCount bool                /* Print count with summary. */
	TagLog       bool                /* Tag logs with argv[0]. */
	PkgType      string              /* Package or Program (default)  */
	Verbose      bool                /* -verbose */
	Imports      map[string]struct{} /* Imported packages. */
}

// SetDefaults makes sure every field of Data has a default value.
func (d *Data) SetDefaults() {
	setDefault(&d.Name, defaultProgramName)
	setDefault(&d.Description, defaultDescription)
	setDefault(&d.Author, defaultAuthorName)
	setDefault(&d.Today, "in the past")
	setDefault(&d.Imports, make(map[string]struct{}))
}

// Clone returns a copy of d.
func (d Data) copy() Data {
	n := d
	n.Imports = maps.Clone(d.Imports)
	return n
}

// CmdDesc gets the command name and description, separated with a ": ".
func (d Data) CmdDesc() string { return d.Name + ": " + d.Description }

// WithImports returns a copy of d with added imports.
func (d Data) WithImports(imports ...string) Data {
	n := d.copy()
	for _, imp := range imports {
		n.Imports[imp] = struct{}{}
	}
	return n
}

// ImportsBlock returns a block of text suitable for use in an imports() block.
// Empty strings will be silently ignored.
func (d Data) ImportsBlock() string {
	var simps, ximps []string
	/* Split the imports into stdlib and external. */
	for imp := range d.Imports {
		if "" == imp {
			continue
		} else if strings.Contains(imp, ".") {
			ximps = append(ximps, imp)
		} else {
			simps = append(simps, imp)
		}
	}
	sort.Strings(simps)
	sort.Strings(ximps)

	/* Buffer for our block. */
	var sb strings.Builder
	sb.WriteString("import (")

	/* Add the two chunks of imports plus spacing. */
	if 0 != len(simps) || 0 != len(ximps) {
		sb.WriteRune('\n')
	}
	for _, imp := range simps {
		fmt.Fprintf(&sb, "\t%q\n", imp)
	}
	if 0 != len(simps) && 0 != len(ximps) {
		/* Blank line between stdlib and external imports. */
		sb.WriteRune('\n')
	}
	for _, imp := range ximps {
		fmt.Fprintf(&sb, "\t%q\n", imp)
	}
	sb.WriteString(")")

	return sb.String()
}

// WithSet returns a copy of d with the named field set to val.  It panics if
// the field cannot be found or if val is not an appropriate type for the
// field.
func (d Data) WithSet(field string, val any) Data {
	ret := d.copy()

	/* Set the field. */
	f := reflect.ValueOf(&ret).Elem().FieldByName(field)
	if !f.IsValid() {
		panic(fmt.Sprintf("Data has no field named %s", field))
	}
	f.Set(reflect.ValueOf(val))

	return ret
}
