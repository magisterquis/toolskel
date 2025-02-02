#!/bin/sh
#
# template.t
# Make sure templates work as expected
# By J. Stuart McMurray
# Created 20250201
# Last Modified 20250201

set -e

. ./t/shmore.subr

tap_plan 3

# Make sure we're not using html/template anywhere.
GOT="$(find . -type f \
        \! -name '*.swp' \
        \! -path './t/*' \
        -exec grep html/template {} + ||:)"
tap_is "$GOT" "" "Not using html/template" "$0" $LINENO

# Make sure template-generated .go files with package main run.  We can use
# the test data for this.
subtest() {
        for FN in $(find . -path '*/testdata/*.go'); do # It'll be ok.
                # Only care about runnable programs.
                if ! egrep -q '^package main$' "$FN"; then
                        continue
                fi
                go run "$FN" >/dev/null
                tap_ok "$?" "Go program ran ok - $FN" "$0" $LINENO
        done
}
tap_subtest "Go template programs run ok" subtest "$0" $LINENO

# Make sure template-generated makefiles parse ok.  We can use the test data
# for this.
subtest() {
        for FN in $(
                find . \
                        -path '*/testdata/*.mk' -o \
                        -path '*/testdata/*Makefile*'
        ); do # It'll be ok.
                make -p "$FN" >/dev/null
                tap_ok "$?" "Makefile parsed ok - $FN" "$0" $LINENO
        done
}
tap_subtest "Template Makefiles parse ok" subtest "$0" $LINENO

# vim: ft=sh
