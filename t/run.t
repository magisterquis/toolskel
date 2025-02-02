#!/bin/sh
#
# files.t
# Test file generation
# By J. Stuart McMurray
# Created 20250201
# Last Modified 20250201

set -e

. ./t/shmore.subr

tap_plan 3

# Will we get a -h happily?
go run . -h 2>/dev/null
tap_ok $? "Ran with -h ok" "$0" $LINENO

# Will we get a warning if we don't specify a file?
GOT=$(! go run . 2>&1)
tap_ok "$?" "(Negated) Non-zero exit with no file to generate" "$0" $LINENO
tap_like \
        "$GOT" \
        '(?m)^\d{4}/\d\d/\d\d \d\d:\d\d:\d\d Need something to generate
exit status 1$' \
        "Warning message with no file to generate" \
        "$0" $LINENO

# vim: ft=sh
