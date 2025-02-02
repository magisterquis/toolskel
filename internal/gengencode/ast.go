package main

/*
 * ast.go
 * Parse the skeleton into an AST.
 * By J. Stuart McMurray
 * Created 20250130
 * Last Modified 20250130
 */

import (
	_ "embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"time"
)

// skel is a source file skeleton, for easy AST creation.
//
//go:embed _skel.go
var skel string

// newAST returns a new skeleton AST from _skel.go consisting of a generated
// comment, package declaration, and imports block.
// The date is speciifed to make for deterministic tests.
func newAST(date time.Time) (*token.FileSet, *ast.File) {
	/* Generate an AST into which to stick methods. */
	fset, f := parseSkel(nil, date)

	/* Remove the template method. */
	f.Decls = f.Decls[:0]

	return fset, f
}

// templateMethod returns a method declaration which should be tweaked and
// added to an ast.File.
// The date is speciifed to make for deterministic tests.
func templateMethod(fset *token.FileSet, date time.Time) *ast.FuncDecl {
	/* Get a new copy of our skeleton. */
	fset, f := parseSkel(fset, date)
	/* Grab the function out of it. */
	sf := f.Decls[0]
	/* Remove the file from the fileset, to save memory. */
	fset.RemoveFile(fset.File(sf.Pos()))

	return sf.(*ast.FuncDecl)
}

// parseSkel parses skel into an AST and makes sure it has a method declaration
// and nothing else.
// If fset is not nil it is used and returned.  Otherwise a new tonen.FileSet
// is created and returned.
// The date is speciifed to make for deterministic tests.
func parseSkel(
	fset *token.FileSet,
	date time.Time,
) (*token.FileSet, *ast.File) {
	/* Parse the skel. */
	if nil == fset {
		fset = token.NewFileSet()
	}
	f, err := parser.ParseFile(
		fset,
		"",
		fmt.Sprintf(skel, date.Format("20060102")),
		parser.ParseComments|parser.SkipObjectResolution,
	)
	if nil != err { /* Should never happen. */
		panic(fmt.Sprintf("generating AST: %s", err))
	} else if 1 != len(f.Decls) {
		panic("AST didn't have exactly one declaration")
	}

	/* Make sure our one declaration is a method. */
	if _, ok := f.Decls[0].(*ast.FuncDecl); !ok {
		panic(fmt.Sprintf("second decl %T, not a method", f.Decls[1]))
	}

	return fset, f
}
