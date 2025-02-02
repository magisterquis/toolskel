Code Generator
==============
Generates progras

Adding Templates
----------------
To add a new template, add a `foo.tmpl` in this directory, re-run
`go generate`, and add tests under all the directories in
[`testdata/TestParams_Generation`](./testdata/TestParams_Generation).  These
should be caught by `make test` if missing.

Template Functions
------------------
Aside from the normal template functions provided by
[`text/template`](https://pkg.go.dev/text/template), the following are
available:

Name         | args             | Description
-------------|------------------|------------
`re_replace` | `re`, `rep`, `s` | Replaces all substrings of `s` which match the regex `re` with the string `rep`
