#!/bin/ksh
#
# files.t
# Make sure our Makefile is the same as what's generated
# By J. Stuart McMurray
# Created 20250222
# Last Modified 20250310

set -euo pipefail

. ./t/shmore.subr

tap_plan 7

# Get the sort of Makefile we expect.
TMPD=$(mktemp -td)
trap 'rm -rf ${TMPD}; tap_done_testing' EXIT
go run . -quiet -dir "$TMPD" -makefile
tap_ok "$?" "Program ran" "$0" $LINENO
NEW="$TMPD/Makefile"
[ -f "$NEW" ]
tap_ok "$?" "Makefile generated" "$0" $LINENO

# And put our own Makefile there for comparison.
CURRENT="$TMPD/Current"
cp Makefile "$CURRENT"

# At this point, commands failing won't bork future tests.
set +e

# For both, set the Created/Last Modified dates to known values
for F in $NEW $CURRENT; do
        sed -Ei 's/^# (Created|Last Modified) [[:digit:]]{8}/# \1 12345678/' $F
        # Make sure it worked
        for l in Created "Last Modified"; do
                GOT=$(egrep "^# $l" $F)
                WANT="# $l 12345678"
                tap_is "$GOT" "$WANT" "$l date correct - $F" "$0" $LINENO
        done
done

# Make sure it's close enough to ours.
GOT="$(diff "$CURRENT" "$NEW")"
tap_is "$GOT" "" "Current and generated Makefiles the same" "$0" $LINENO

# vim: ft=sh
