#!/bin/ksh
#
# txtar.t
# Test txtar generation
# By J. Stuart McMurray
# Created 20250201
# Last Modified 20250407

set -uo pipefail

. ./t/shmore.subr

tap_plan 12

TESTDATA=t/testdata/txtar

# Update all the files with macros
make -C "$TESTDATA" -f ../../updatemacros.mk -s
tap_ok $? "Macro files up-to-date" "$0" $LINENO

# test_gen tests creation of $1.
test_gen() {
        tap_plan 3
        # Generate the archive
        GOT=$(go run . \
                "-$1" \
                -author "test_author" \
                -name "" \
                -quiet \
                -today 25251225 \
                -txtar)
        tap_ok $? "Ran successfully" "$0" $LINENO

        # Remove the generated date
        L1=$(/bin/echo "$GOT" | head -n 1)
        tap_like \
                "$L1" \
                '^Created by toolskel '\
'((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)$' \
                "Comment correct" \
                "$0" $LINENO

        # Work out if it's different.
        WANTF="t/testdata/txtar/$1"
        DIFF=$(/bin/echo "$GOT" | tail -n +2 | diff -u - "$WANTF" ||:)
        tap_is "$DIFF" "" "Output correct" "$0" $LINENO
}

# Test ALL the things
for FLAG in t/testdata/txtar/*; do
        # Don't build M4 macro files or the Makefile
        if [[ "$FLAG" == *.m4 ]]; then
                continue
        fi
        FLAG=${FLAG##*/}
        subtest() { test_gen "$FLAG"; }
        tap_subtest "Generate -$FLAG" "subtest" "$0" $LINENO
done

# vim: ft=sh
