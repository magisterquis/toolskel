#!/bin/ksh
#
# using_basic_tests.t
# Make sure our basic_tests.t is the same as what's generated
# By J. Stuart McMurray
# Created 20250310
# Last Modified 20250316

set -uo pipefail

. ./t/shmore.subr

tap_plan 3

# Filenames
TBTF="t/basic_tests.t"
GENF="$(mktemp -t gen.tmp.XXXXXXXXXX)"
CURF="$(mktemp -t cur.tmp.XXXXXXXXXX)"
AWKP="NR < 5 || 7 < NR"
trap 'rm $GENF $CURF; tap_done_testing' EXIT

# Get the sort of basic_tests.t we expect.
go run . -quiet -txtar -basic-tests | tail -n +3 | awk "$AWKP" >$GENF
tap_ok "$?" "Generated interpolatedless $TBTF in $GENF" "$0" $LINENO

# And get ours, less the header.
awk "$AWKP" "$TBTF" >$CURF
tap_ok "$?" "Put nameless/dateless $TBTF in $CURF" "$0" $LINENO

# Make sure it's close enough to ours.
GOT="$(diff -u "$GENF" "$CURF" ||:)"
tap_is "$GOT" "" "Current and generated $TBTF the same" "$0" $LINENO

# vim: ft=sh
