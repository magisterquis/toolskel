package gencode

/*
 * templatefuncs.go
 * Functions available in templates
 * By J. Stuart McMurray
 * Created 20250201
 * Last Modified 20250201
 */

import (
	"fmt"
	"regexp"
	"text/template"
)

// Names of functions in newFuncMapc's returned map.
const (
	fnREReplace = "re_replace"
)

// newFuncMap returns a new function map with our custom template functions.
func newFuncMap() template.FuncMap {
	return template.FuncMap{
		fnREReplace: replaceAll,
	}
}

// replaceAll replaces all substrings of s which match the regex r with rep.
func replaceAll(r, rep, s string) (string, error) {
	re, err := regexp.Compile(r)
	if nil != err {
		return "", fmt.Errorf("compiling %s: %w", r, err)
	}
	return re.ReplaceAllString(s, rep), nil
}
