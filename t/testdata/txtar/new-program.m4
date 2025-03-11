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

.PHONY: all build test gotest provetest help install clean

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
			{ print "Long usage line: " $0; exit 1 }\
	'

provetest: ## Run tests with prove(1) if ./t exists
.if exists(./t/)
	prove -It --directives
.endif

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
#!/bin/sh
#
# debug.t
# Make sure we don't have debugging things left in
# By test_author
# Created 25251225
# Last Modified 25251225

set -eu
if (set -o pipefail 2>/dev/null); then set -o pipefail; fi

. t/shmore.subr

tap_plan 3

# Make sure we're not using MQD.
GOT="$(go run . -h </dev/null 2>&1 | egrep 'MQD DEBUG PACKAGE LOADED$' ||:)"
tap_is "$GOT" "" "Not using github.com/magisterquis/mqd" "$0" $LINENO

# Make sure we didn't leave any stray DEBUGs lying about.
GOT="$(egrep -InR '#[[:space:]]*DEBUG' | sort -u ||:)"
tap_is \
        "$GOT" \
        "" \
        "No files with DEBUG comments" \
        "$0" $LINENO

# Should get happy help output.  Not exiting cleanly is ok, and part of
# testing.
set +e
go run . -h 2>/dev/null
tap_ok $? "Ran with -h ok" "$0" $LINENO
set -e

# vim: ft=sh
-- t/shmore.subr --
m4_paste(m4_shmore)m4_dnl
