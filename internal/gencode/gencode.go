// Package gencode generates code for toolskel.  It's where everything happens.
package gencode

/*
 * gencode.go
 * Generate toolskel code.
 * By J. Stuart McMurray
 * Created 20230425
 * Last Modified 20230425
 */

import (
	"fmt"
	"io"
	"reflect"
	"text/template"
)

// DefaultTType is the default template to use.
const DefaultTType = "simple"

// Templates are the parsed templates
var templates = make(map[string]*template.Template)

func init() {
	mustParseTemplates()
}

// Generate does the code generation itself.
func Generate(w io.Writer, tType string, data Data) error {
	/* Make sure we have a template type. */
	setDefault(&tType, DefaultTType)

	/* Make sure all of the fields are filled. */
	data.SetDefaults()

	/* Get the template for this type, making sure we've parsed the
	templates. */
	tmpl, ok := templates[tType]
	if !ok {
		return fmt.Errorf("unknown tool type %q", tType)
	}

	/* Emit boilerplate. */
	return tmpl.Execute(w, data)
}

// setDefault sets *p to T if *p is the zero value for its type.  If p is nil,
// SetDefault panics.
func setDefault[T string | map[string]struct{}](p *T, def T) {
	/* Doesn't work with nil. */
	if nil == p {
		panic("setDefault: nil pointer")
	}

	/* If we've got the zero value, set the default. */
	if reflect.ValueOf(p).Elem().IsZero() {
		*p = def
	}
}
