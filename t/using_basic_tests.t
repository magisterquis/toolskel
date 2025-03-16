#!/bin/ksh
#
# using_basic_tests.t
# Make sure our basic_tests.t is the same as what's generated
# By J. Stuart McMurray
# Created 20250310
# Last Modified 20250316

set -euo pipefail

. ./t/shmore.subr

tap_plan 3

# Filenames
TESTF="t/basic_tests.t"
WANTF="$(mktemp -t want.tmp.XXXXXXXXXX)"
HAVEF="$(mktemp -t have.tmp.XXXXXXXXXX)"
trap 'rm $WANTF $HAVEF; tap_done_testing' EXIT

# Get the sort of basic_tests.t we expect.
go run . -quiet -txtar -basic-tests | tail -n +10 >$WANTF
tap_ok "$?" "Generated headerles $TESTF in $WANTF" "$0" $LINENO

# And get ours, less the header.
tail -n +8 "$TESTF" >$HAVEF
tap_ok "$?" "Put headerless $TESTF in $HAVEF" "$0" $LINENO

# Make sure it's close enough to ours.
GOT="$(diff -u "$HAVEF" "$WANTF" ||:)"
tap_is "$GOT" "" "Current and generated $TESTF the same" "$0" $LINENO

# vim: ft=sh
