#!/bin/ksh
#
# files.t
# Test file generation
# By J. Stuart McMurray
# Created 20250201
# Last Modified 20250311

set -euo pipefail

. ./t/shmore.subr

tap_plan 2

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
