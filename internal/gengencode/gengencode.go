// Program gengencode generates code generators for gencode.
package main

/*
 * gengencode.go
 * Generate code generators for gencode.
 * By J. Stuart McMurray
 * Created 20250130
 * Last Modified 20250316
 */

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	// tmplSuffix is the suffix for template files.
	tmplSuffix = `.tmpl`
)

// errNoTemplates is returned by do when a directory has no template files.
var errNoTemplates = errors.New("no template files found")

// verbose will be turned off unless we get a -verbose
var verbose = log.Printf

func main() {
	var (
		outputFile = flag.String(
			"write",
			"",
			"Output `file`, unset for stdout",
		)
		verbOn = flag.Bool(
			"verbose",
			false,
			"Print method names as they're added",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %s [outputfile]

Generates Go methods on a Params type to call Params.execute for each file
in the current directory named *.tmpl.

See _skel.go.

Options:
`,
			filepath.Base(os.Args[0]),
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* If we're not being verbose, make it a no-op. */
	if !*verbOn {
		verbose = func(string, ...any) {}
	}

	/* Generate the code. */
	b, err := do(".", time.Now())
	if nil != err {
		log.Fatalf("gengencode: Error generating code: %s", err)
	}

	/* If we're just printing the generated code, life's easy. */
	if "" == *outputFile {
		os.Stdout.Write(b)
		return
	}

	/* Work out what we've already got.  If there's no change, no need to
	do anything. */
	old, err := os.ReadFile(*outputFile)
	if nil != err && !errors.Is(err, fs.ErrNotExist) {
		log.Fatalf(
			"gengencode: Error reading %s: %s",
			*outputFile,
			err,
		)
	}
	_, oldBody, _ := strings.Cut(string(old), "\n")
	_, newBody, _ := strings.Cut(string(b), "\n")
	if oldBody == newBody {
		return
	}

	/* Got new data.  Update the file. */
	if err := os.WriteFile(*outputFile, b, 0644); nil != err {
		log.Fatalf(
			"gengencode: Error writing to %s: %s",
			*outputFile,
			err,
		)
	}
}

// do generates Params.Generate* functions for gencode which execute each
// *.tmpl file in the directory where.
// The date is specified to make for deterministic tests.
func do(dir string, date time.Time) ([]byte, error) {
	/* Make a new AST for file-generation. */
	fset, astf := newAST(date)

	/* Template files for which to generate generators. */
	tmplfns, err := filepath.Glob(filepath.Join(dir, "*"+tmplSuffix))
	if nil != err {
		return nil, fmt.Errorf("listing template files: %w", err)
	} else if 0 == len(tmplfns) {
		return nil, errNoTemplates
	}

	/* Make sure the templates parse, and add methods for each one. */
	for _, fn := range tmplfns {
		if err := addMethod(
			fset,
			astf,
			filepath.Base(fn),
			date,
		); nil != err {
			return nil, fmt.Errorf("adding %s: %w", fn, err)
		}
	}

	/* Reformat the source and send it back. */
	b := new(bytes.Buffer)
	if err := format.Node(b, fset, astf); nil != err {
		return nil, fmt.Errorf("formatting generated code: %w", err)
	}
	return b.Bytes(), nil
}

// addMethod adds a method for fn to astf.
func addMethod(
	fset *token.FileSet,
	astf *ast.File,
	fn string,
	date time.Time,
) error {
	/* Method name is the same as the template name, minus the suffix,
	capitalized if need be. */
	mn := cases.Title(language.Und, cases.NoLower).String(
		strings.TrimSuffix(fn, tmplSuffix),
	)

	/* Grab a template for our method. */
	tm := templateMethod(fset, date)

	/* Only things which need updating are the function name and the
	file to be grabbed.  Name's easy. */
	tm.Name = ast.NewIdent(mn)

	/* Filename to be grabbed is slightly complexer, but should be the RHS
	of the first statement in the function.  Bit of a leap of faith here,
	but a panic at this point is probably a better way to tell someone I
	borked the template than a gentle error message. */
	tm.Body.List[0].(*ast.ReturnStmt).
		Results[0].(*ast.CallExpr).
		Args[0].(*ast.BasicLit).
		Value = strconv.QuoteToASCII(fn)

	/* Add the method to the file. */
	astf.Decls = append(astf.Decls, tm)

	verbose("Added method %s", mn)

	return nil

}
