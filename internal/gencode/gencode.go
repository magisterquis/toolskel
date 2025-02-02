// Package gencode generates code from templates.
package gencode

/*
 * gencode.go
 * Generate code from templates.
 * By J. Stuart McMurray
 * Created 20250130
 * Last Modified 20250201
 */

import (
	"bytes"
	"cmp"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

var (
	// tmplFS makes the embedded templates accessible as an fs.FS.
	//
	//go:embed *.tmpl
	tmplFS embed.FS

	tmpl = template.Must(template.New("").Funcs(newFuncMap()).ParseFS(tmplFS, `*.tmpl`))
)

// Defaults for like-named fields in Params.
const (
	DefaultName        = "something"
	DefaultDescription = "A cool Something"
	DefaultAuthor      = "Someone"
)

// Generator is a method on Params which generates a file.
type Generator = func() ([]byte, error)

// Params is the data we use when generating a file.  Its methods generate
// the files suggested by their names, placing the output in the passed-in
// filename.
type Params struct {
	Name        string /* Tool/project/such name. */
	Description string /* One-line description. */
	Author      string /* Author's name. */
	Today       string /* YYYYMMDD. */
}

// ToFile wraps a Generate* function and writes its output to the file named
// fn.  Output is first written to a temporary file which is then atomically
// renamed.  If fn already exists, it will only be overwritten if overwrite is
// true.
func (p Params) ToFile(fn string, generator Generator, overwrite bool) error {
	/* Don't accidentally overwrite fn. */
	fi, err := os.Stat(fn)
	switch {
	case errors.Is(err, fs.ErrNotExist): /* This is the easy case. */
	case nil != err: /* Something went wrong. */
		return fmt.Errorf("getting information about %s: %s", fn, err)
	case !fi.Mode().IsRegular(): /* Probably a directory. */
		return fmt.Errorf("not a regular file")
	case !overwrite: /* Don't clobber. */
		return fmt.Errorf("file already exists")
	}

	/* Execute the template. */
	b, err := generator()
	if nil != err {
		return fmt.Errorf("generating file: %w", err)
	}

	/* Temporary file we'll move over fn, for atomicity. */
	f, err := os.CreateTemp(filepath.Dir(fn), filepath.Base(fn))
	if nil != err {
		return fmt.Errorf("creating temporary file: %w", err)
	}
	defer f.Close()

	/* Add the file to the file. */
	if _, err := f.Write(b); nil != err {
		os.Remove(f.Name()) /* Best effort */
		return fmt.Errorf(
			"writing to temporary file %s: %w",
			f.Name(),
			err,
		)
	}

	/* Atomically update fn. */
	if err := os.Rename(f.Name(), fn); nil != err {
		os.Remove(f.Name()) /* Best effort */
		return fmt.Errorf(
			"renaming temporary file %s to %s: %w",
			f.Name(),
			fn,
			err,
		)
	}

	return nil
}

// execute executes the template named tn.
func (p Params) execute(tn string) ([]byte, error) {
	/* Set defaults. */
	(&p).setDefaults()

	/* Make sure we have the template. */
	t := tmpl.Lookup(tn)
	if nil == t {
		panic(fmt.Sprintf(
			"unknown template %s (need to re-run go generate?)",
			tn,
		))
	}

	/* Execute the template itself. */
	b := new(bytes.Buffer)
	if err := t.Execute(b, p); nil != err {
		return nil, fmt.Errorf("executing template: %w", err)
	}

	return b.Bytes(), nil
}

// Make sure a params has default values.
func (p *Params) setDefaults() {
	setDef := func(p *string, def string) { *p = cmp.Or(*p, def) }
	for p, def := range map[*string]string{
		&p.Name:        DefaultName,
		&p.Description: DefaultDescription,
		&p.Author:      DefaultAuthor,
	} {
		setDef(p, def)
	}
	setDef(&p.Today, time.Now().Format("20060102"))
}

//go:generate go run ../gengencode -verbose -write tmpl.go
