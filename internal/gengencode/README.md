Generator for Code Generator
============================
Program to generate the methods on [gencode.Params](../gencode/gencode.go)
which create files from [gencode's templates](../gencode).

To use, be in [gencode](../gencode) and call it.  Should be called from a
`//go:generate` line, though.

It should probably just use templates and not do
[AST](https://pkg.go.dev/go/ast) things.
