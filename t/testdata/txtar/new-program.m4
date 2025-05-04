m4_changequote(dummy)m4_dnl
-- something.go --
// Program something - A cool Something
package main

/*
 * something.go
 * A cool Something
 * By test_author
 * Created 25251225
 * Last Modified 25251225
 */

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() { os.Exit(rmain()) }
func rmain() int {
	/* Command-line flags. */
	var (
	/* TODO: Add flags. */
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %s [options]

A cool Something

Options:
`,
			filepath.Base(os.Args[0]),
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* TODO: Meat and Potatoes. */

	return 0
}
-- Makefile --
# Makefile
# Build something
# By test_author
# Created 25251225
# Last Modified 25251225

BINNAME       != basename $$(pwd)
GOBUILDFLAGS   = -trimpath -ldflags "-w -s"
GOTESTFLAGS   += -timeout 3s
SHMORESUBR     = t/shmore.subr
SHMOREURL      = https://raw.githubusercontent.com/magisterquis/shmore/refs/heads/master/shmore.subr

.PHONY: all build test gotest provetest help install update clean

all: test build ## Build ALL the things (default)

${BINNAME}:
	go build ${GOBUILDFLAGS} -o ${BINNAME}

build: ${BINNAME}

test: gotest provetest ## Run ALL the tests

gotest: ## Run go-specific tests
	go test ${GOBUILDFLAGS} ${GOTESTFLAGS} ./...
	go vet ${GOBUILDFLAGS} ./...
	staticcheck ./...
	go run ${GOBUILDFLAGS} . -h 2>&1 |\
	awk '\
		/^Options:$$|MQD DEBUG PACKAGE LOADED$$/\
			{ exit }\
		/^Usage: /\
			{ sub(/^Usage: [^[:space:]]+\//, "Usage: ") }\
		/.{80,}/\
			{ print "Long usage line: " $$0; exit 1 }\
	'

provetest: ## Run tests with prove(1) if ./t exists
.if exists(./t/)
	prove -It --directives
.endif

update: ## Fetch the latest Shmore and up-to-date Go things
	curl\
		--fail\
		--show-error\
		--silent\
		--output ${SHMORESUBR}.new\
		${SHMOREURL}
	diff -q ${SHMORESUBR} ${SHMORESUBR}.new >/dev/null &&\
		rm ${SHMORESUBR}.new ||\
		mv ${SHMORESUBR}.new ${SHMORESUBR}
	go get -t -u go ./...
	go mod tidy

install: ## Install to GOBIN ($GOPATH/bin or $HOME/go/bin)
	go install ${GOBUILDFLAGS}

clean: ## Remove built things
	rm -rf ${BINNAME}

help: .NOTMAIN ## This help
	@perl -ne '/^(\S+?):+.*?##\s*(.*)/&&print"$$1\t-\t$$2\n"' \
		${MAKEFILE_LIST} | column -ts "$$(printf "\t")"
-- staticcheck.conf --
checks = ["all", "-ST1017"]
-- README.md --
something
=========
A cool Something

Quickstart
----------
1. Write a quickstart...

Usage
-----
```
TODO: Usage
```
-- .gitignore --
something
something-*-*
.*.swp
-- t/basic_tests.t --
#!/bin/ksh
#
# basic_tests.t
# Make sure our code is up-to-date and doesn't have debug things.
# By test_author
# Created 25251225
# Last Modified 25251225

set -uo pipefail

. t/shmore.subr

NTEST=6
tap_plan "$NTEST"

# Make sure we didn't leave any stray DEBUGs or TAP_TODOs lying about.
GOT=$(egrep -InR '(#|\*)[[:space:]]*()DEBUG' | sort -u)
tap_is "$GOT" "" "No files with DEBUG comments" "$0" $LINENO
GOT=$(egrep -In  'TAP_TODO[=]' t/*.t | sort -u)
tap_is "$GOT" "" "No TAP_TODO's" "$0" $LINENO

# These checks assume we're writing a Go program.
if [[ -f ./go.mod ]]; then
        # TMPD is where we'll put our temporary program
        TMPD=$(mktemp -td)
        trap 'rm -rf ${TMPD}; tap_done_testing' EXIT

        # Make sure we're not using MQD.
        GOT="$(go run . -h </dev/null 2>&1 |
                egrep 'MQD DEBUG PACKAGE LOADED$')"
        tap_is "$GOT" "" "Not using github.com/magisterquis/mqd" "$0" $LINENO

        # Should get happy help output.  We can't use go run here because it
        # doesn't properly propagate the exit status.
        go build -o "$TMPD/tb"
        "$TMPD/tb" -h 2>/dev/null
        tap_is $? 0 "Running with -h exits happily" "$0" $LINENO

        # Make sure we don't need to update anything.
        GOT="$(go list \
                -u \
                -f '{{if (and (not (or .Main .Indirect)) .Update)}}
                        {{- .Path}}: {{.Version}} -> {{.Update.Version -}}
                {{end}}' \
                -m all)"
        tap_is "$GOT" "" "Packages up-to-date" "$0" $LINENO
        # Idea stolen from https://github.com/fogfish/go-check-updates

        # Make sure we're using the latest Go as well.
        GOT="$(go list \
                -u \
                -f '{{if (and .Update .Update.Version) -}}
                        go {{.Version}} -> {{.Update.Version}}
                {{- end}}' \
                -m go)"
        tap_is "$GOT" "" "Latest Go version will be used" "$0" $LINENO
else
        tap_skip "Not a Go program" $((NTEST-2))
fi

# vim: ft=sh
-- t/shmore.subr --
m4_paste(m4_shmore)m4_dnl
