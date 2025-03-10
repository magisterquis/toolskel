#!/bin/sh
#
# txtar.t
# Test txtar generation
# By J. Stuart McMurray
# Created 20250201
# Last Modified 20250310

set -euo pipefail

. ./t/shmore.subr

tap_plan 9

# test_gen tests creation of $1.
test_gen() {
        tap_plan 4
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
        GOT=$(/bin/echo "$GOT" | tail -n +2)

        # Work out what it should be
        WANTF="t/testdata/txtar/$1"
        WANT="$(cat "$WANTF")"
        tap_isnt "$WANT" "" "Read WANT from $WANTF" "$0" $LINENO

        # Make sure it is what it should be
        tap_is "$GOT" "$WANT" "Output correct" "$0" $LINENO
}

# Test ALL the things
for FLAG in t/testdata/txtar/*; do
        FLAG=${FLAG##*/}
        subtest() { test_gen "$FLAG"; }
        tap_subtest "Generate -$FLAG" "subtest" "$0" $LINENO
done

# vim: ft=sh
