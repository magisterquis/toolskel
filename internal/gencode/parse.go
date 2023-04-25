package gencode

/*
 * parse.go
 * Parse templates
 * By J. Stuart McMurray
 * Created 20230421
 * Last Modified 20230425
 */

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
)

// templateSuffix is the suffix on each template.
const templateSuffix = ".tmpl"

var (
	// baseTemplate holds the base template.
	//
	//go:embed base.tmpl
	baseTemplate string

	// templateFS holds the underlying template files
	//
	//go:embed templates/*.tmpl
	templateFS embed.FS
)

// mustParseTemplates parses the templates into Templates.  The base is taken
// from baseTemplate, and the templates in templateFS are used to populate
// Templates.  MustParseTemplates panics on error.
func mustParseTemplates() {
	/* Get the base template. */
	baseT := template.Must(template.New("base").Parse(baseTemplate))

	/* Work out the other types of templates we have. */
	if err := fs.WalkDir(templateFS, ".", func(
		path string,
		d fs.DirEntry,
		err error,
	) error {
		/* Errors are unpossible. */
		if nil != err {
			return fmt.Errorf("walking %q: %w", path, err)
		}

		/* Can't really use directories. */
		if d.IsDir() {
			return nil
		}

		/* We shouldn't have anything not ending in .tmpl. */
		if !strings.HasSuffix(path, templateSuffix) {
			return fmt.Errorf("unexpected filename: %q", path)
		}

		/* Template name should be unique. */
		tn := filepath.Base(strings.TrimSuffix(path, templateSuffix))
		if _, ok := templates[tn]; ok {
			return fmt.Errorf("template already defined: %q", tn)
		}

		/* Get the template's body. */
		b, err := templateFS.ReadFile(path)
		if nil != err {
			return fmt.Errorf(
				"getting embeddede %q: %w",
				path,
				err,
			)
		}

		/* Parse into a template. */
		templates[tn], err = template.Must(baseT.Clone()).
			Parse(string(b))
		if nil != err {
			return fmt.Errorf("parsing %q: %w", path, err)
		}

		return nil
	}); nil != err {
		panic(fmt.Sprintf("parsing templates: %s", err))
	}
}
