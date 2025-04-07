#!/bin/ksh
#
# go_updates.t
# Make sure the Checker of Go Updates works
# By J. Stuart McMurray
# Created 20250407
# Last Modified 20250407

set -uo pipefail

. ./t/shmore.subr

tap_plan 5

set -e # This all needs to work
GOTF=$(mktemp -t) # Command output
TD=$(mktemp -d)   # Project directory
trap 'R=$?; rm -rf "$GOTF" "$TD"; (exit $R); tap_done_testing' EXIT
# Set up a go project with old dependencies
cp -r t/testdata/go_updates/old_versions/* $TD
mkdir "$TD/t"
cp t/shmore.subr "$TD/t/"
tap_pass "Put go project in $TD" "$0" $LINENO
set +e
WANTF=$(mktemp -t) # What we want

# diff_is makes sure $1 and stdin (i.e. want) are the same, but with diff
#
# Options:
# $1 - The got
# $2 - Test name
# $3 - Filename
# $4 - Line number
diff_is() {
        echo "$1" >$GOTF
        GOT=$(diff -u "/dev/stdin" "$GOTF")
        tap_is "$GOT" "" "$2" "$3" "$4"
}

# Add our magic makefile and test
go run . -dir "$TD" -quiet -makefile -basic-tests 2>&1 |
        awk '! /^go: downloading/'
RET=$?
tap_ok $? "Generated makefile and basic tests" "$0" $LINENO
if [[ 0 != $RET ]]; then
        exit 4
fi

# Check for out-of-date things
GOT=$(cd "$TD" &&
        HARNESS_ACTIVE=1 ksh ./t/basic_tests.t 2>&1 |
        sed -E 's/ -> [^[:space:]]+/ -> XXX/')
diff_is "$GOT" "Old versions detected" "$0" $LINENO <<'_eof'
1..6
ok 1 - No files with DEBUG comments
ok 2 - No TAP_TODO's
ok 3 - Not using github.com/magisterquis/mqd
ok 4 - Running with -h exits happily
not ok 5 - Packages up-to-date

#   Failed test 'Packages up-to-date'
#   at ./t/basic_tests.t line 46.
#          got: 'golang.org/x/net: v0.37.0 -> XXX
# golang.org/x/text: v0.23.0 -> XXX
#     expected: ''
not ok 6 - Latest Go version will be used

#   Failed test 'Latest Go version will be used'
#   at ./t/basic_tests.t line 56.
#          got: 'go 1.24.1 -> XXX
#     expected: ''
# Looks like you failed 2 tests of 6.
_eof

# Update ALL the things.
GOT=$(make -s -C $TD update 2>&1 |
        awk '! /^go: downloading/' |
        sed -E 's/ => [^[:space:]]+/ => XXX/')
diff_is "$GOT" "Updated happily" "$0" "$LINENO" <<'_eof'
go: upgraded go 1.24.1 => XXX
go: upgraded golang.org/x/net v0.37.0 => XXX
go: upgraded golang.org/x/text v0.23.0 => XXX
_eof

# Make sure the updates worked.
GOT=$(cd "$TD" && ksh ./t/basic_tests.t 2>&1)
diff_is "$GOT" "Everything up-to-date" "$0" $LINENO <<'_eof'
1..6
ok 1 - No files with DEBUG comments
ok 2 - No TAP_TODO's
ok 3 - Not using github.com/magisterquis/mqd
ok 4 - Running with -h exits happily
ok 5 - Packages up-to-date
ok 6 - Latest Go version will be used
_eof

# vim: ft=sh
