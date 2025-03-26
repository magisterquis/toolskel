#!/bin/sh
#
# debug.t
# Make sure we don't have debugging things left in
# By test_author
# Created 25251225
# Last Modified 25251225

set -u
if (set -o pipefail 2>/dev/null); then set -o pipefail; fi

. t/shmore.subr

tap_plan 3

# TMPD is where we'll put our temporary program
TMPD=$(mktemp -td)
trap 'rm -rf ${TMPD}; tap_done_testing' EXIT

# Make sure we're not using MQD.
GOT="$(go run . -h </dev/null 2>&1 | egrep 'MQD DEBUG PACKAGE LOADED$')"
tap_is "$GOT" "" "Not using github.com/magisterquis/mqd" "$0" $LINENO

# Make sure we didn't leave any stray DEBUGs lying about.
GOT="$(egrep -InR '#[[:space:]]*DEBUG' | sort -u)"
tap_is \
        "$GOT" \
        "" \
        "No files with DEBUG comments" \
        "$0" $LINENO

# Should get happy help output.  We can't use go run here because it doesn't
# properly propagate the exit status.
go build -o "$TMPD/tb"
"$TMPD/tb" -h 2>/dev/null
tap_is $? 0 "Running with -h exits happily" "$0" $LINENO

# vim: ft=sh
