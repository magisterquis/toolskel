package gencode

/*
 * list_test.go
 * Tests for list.go
 * By J. Stuart McMurray
 * Created 20230421
 * Last Modified 20230425
 */

import (
	"bytes"
	"testing"
)

// TestTemplateDescriptions makes sure all templates have descriptions.
func TestTemplateDescriptions(t *testing.T) {
	for tn, tmpl := range templates {
		tmpl := tmpl /* :C */
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			dt := tmpl.Lookup(descriptionTemplate)
			if nil == dt {
				t.Errorf("No description subtemplate")
				return
			}
			var buf bytes.Buffer
			if err := dt.Execute(&buf, nil); nil != err {
				t.Errorf(
					"Error extracting description: %s",
					err,
				)
				return
			}
			if 0 == buf.Len() {
				t.Errorf("Empty description")
				return
			}
		})
	}
}
