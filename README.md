Tool Skeleton Generator
=======================
Generates boilerplate for small tools written in Go.

Expected userbase size: 1

Tool Types
----------
The currently-available tool types are

Type       | Description
-----------|------------
`library`  | Just headers, for a library
`parallel` | Parallel task executor
`simple `  | A no-frills tool

Usage
-----
```
Usage: toolskel [options] [toolname [tool description...]]

Generates boilerplate for a tool written in Go.

Options:
  -author name
    	Author's name (default "Stuart McMurray")
  -list-types
    	List available tool types
  -no-date
    	Do not set the Created/Modified date
  -summary-count
    	Generated code's summary prints a completed task count
  -tag-log
    	Tag log output with argv[0]
  -type type
    	Tool type (see -list-types) (default "simple")
```

Quickstart
----------
```sh
go install github.com/magisterquis/toolskel@latest
toolskel -h
toolskel -author 'Darth Vader' findrebels Finds rebel scum > tool.go
vi ./tool.go
```

Building and Testing
--------------------
In most cases, `go install` should be sufficient.  The [Makefile](./Makefile)
is intended for use during development.  The Makefile assumes that
[`goimports`](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) and 
[`Staticcheck`](https://staticcheck.io) are available.

Adding Templates
----------------
1.  Add a template to
    [`internal/gencode/templates`](./internal/gencode/templates) which should
    replace blocks in
    [`internal/gencode/base.tmpl`](./internal/gencode/base.tmpl).  Make sure to
    update the `description` block.
2.  Add a testcase or three to
    [`internal/gencode/toolskel_test.go`](.internal/gencode/toolskel_test.go).
3.  Generate a test copy of the output with something like
    ```sh
go run . -author '' -no-date -type $NEWTYPE >> internal/gencode/tests/newtype.go
```
    The name should be the same as Testcases[yours].name, with `s,/,_,g`.
4.  Run the tests with
    ```sh
make tests
```
