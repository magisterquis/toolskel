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


# Set up a go project with old dependencies
set -e # This all needs to work
TD=$(mktemp -d)
trap 'rm -rf "$TD"; tap_done_testing' EXIT # Make sure it's removed
cp -r t/testdata/go_updates/old_versions/* $TD
mkdir "$TD/t"
cp t/shmore.subr "$TD/t/"
tap_pass "Put go project in $TD" "$0" $LINENO
set +e

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
WANT=$(cat <<'_eof'
1..5
ok 1 - No files with DEBUG comments
ok 2 - Not using github.com/magisterquis/mqd
ok 3 - Running with -h exits happily
not ok 4 - Packages up-to-date

#   Failed test 'Packages up-to-date'
#   at ./t/basic_tests.t line 47.
#          got: 'golang.org/x/net: v0.37.0 -> XXX
# golang.org/x/text: v0.23.0 -> XXX
#     expected: ''
not ok 5 - Latest Go version will be used

#   Failed test 'Latest Go version will be used'
#   at ./t/basic_tests.t line 57.
#          got: 'go 1.24.1 -> XXX
#     expected: ''
# Looks like you failed 2 tests of 5.
_eof
)
tap_is "$GOT" "$WANT" "Old versions detected" "$0" $LINENO

# Update ALL the things.
GOT=$(make -s -C $TD update 2>&1 |
        awk '! /^go: downloading/' |
        sed -E 's/ => [^[:space:]]+/ => XXX/')
WANT=$(cat <<'_eof'
go: upgraded go 1.24.1 => XXX
go: upgraded golang.org/x/net v0.37.0 => XXX
go: upgraded golang.org/x/text v0.23.0 => XXX
_eof
)
tap_is "$GOT" "$WANT" "Updated happily" "$0" "$LINENO"

# Make sure the updates worked.
GOT=$(cd "$TD" && ksh ./t/basic_tests.t 2>&1)
WANT=$(cat <<'_eof'
1..5
ok 1 - No files with DEBUG comments
ok 2 - Not using github.com/magisterquis/mqd
ok 3 - Running with -h exits happily
ok 4 - Packages up-to-date
ok 5 - Latest Go version will be used
_eof
)

tap_is "$GOT" "$WANT" "Everything up-to-date" "$0" $LINENO

# vim: ft=sh
