#!/bin/ksh
#
# go_updates.t
# Make sure the Checker of Go Updates works
# By J. Stuart McMurray
# Created 20250407
# Last Modified 20260511

set -euo pipefail

. ./t/shmore.subr

tap_plan 5

# diff_is makes sure $1 and stdin (i.e. want) are the same, but with diff
#
# Options:
# $1 - The got
# $2 - Test name
# $3 - Filename ($0)
# $4 - Line number
diff_is() {
        local _got=$1 _testname=$2 _file=$3 _line=$4
        GOT=$(diff -u \
                -L "Expected output" \
                -L "Actual output" \
                /dev/stdin /dev/fd/3 3<<_eof ||:
$_got
_eof
)
        tap_is "$GOT" "" "$_testname" "$_file" "$_line"
}

# Set up a go project with old dependencies
TD=$(mktemp -d)
trap 'rm -rf "$TD"; tap_done_testing' EXIT # Make sure it's removed
cp -r t/testdata/go_updates/old_versions/* $TD
mkdir "$TD/t"
cp t/shmore.subr "$TD/t/"
tap_pass "Put go project in $TD" "$0" $LINENO

# Add our magic makefile and test
go run . -dir "$TD" -quiet -makefile -basic-tests 2>&1 |
        awk '! /^go: downloading/'
tap_pass "Generated makefile and basic tests" "$0" $LINENO

# Check for out-of-date things
GOT=$(cd "$TD" &&
        { HARNESS_ACTIVE=1 ksh ./t/basic_tests.t ||: ; } 2>&1 |
        sed -E 's/ -> [^[:space:]]+/ -> XXX/')
diff_is "$GOT" "Old versions detected" "$0" $LINENO <<'_eof'
1..7
ok 1 - No files with unexpected DEBUG comments
ok 2 - No files with unexpected TODO comments
ok 3 - No TAP_TODO's
ok 4 - Not using github.com/magisterquis/mqd
ok 5 - Running with -h exits happily
not ok 6 - Packages up-to-date

#   Failed test 'Packages up-to-date'
#   at ./t/basic_tests.t line 63.
#          got: 'golang.org/x/net: v0.37.0 -> XXX
# golang.org/x/text: v0.23.0 -> XXX
#     expected: ''
not ok 7 - Latest Go version will be used

#   Failed test 'Latest Go version will be used'
#   at ./t/basic_tests.t line 73.
#          got: 'go 1.24.1 -> XXX
#     expected: ''
# Looks like you failed 2 tests of 7.
_eof

# Update ALL the things.
BIN_VERSION=$(cd $TD && go version | awk '{print $3}')
GOT=$(make -s -C $TD update 2>&1 |
        awk '! /^go: downloading/' |
        sed -E 's/ => [^[:space:]]+/ => XXX/')
WANT=$(cat <<_eof
go: upgraded go 1.24.1 => XXX
go: upgraded golang.org/x/net v0.37.0 => XXX
go: upgraded golang.org/x/text v0.23.0 => XXX
_eof
)
# If we have a different compiler version from the module version, we'll expect
# to also get a message that Go is switching compiler versions.
MOD_VERSION="go$(cd $TD && go list -f {{.GoVersion}} -m)"
if [[ "$BIN_VERSION" != "$MOD_VERSION" ]]; then
        WANT="go: updating go.mod requires go >= ${MOD_VERSION#go}; \
switching to $MOD_VERSION
$WANT"
fi
diff_is "$GOT" "Updated correctly" "$0" "$LINENO" <<_eof
$WANT
_eof

# Make sure the updates worked.
GOT=$(cd "$TD" && ksh ./t/basic_tests.t 2>&1)
diff_is "$GOT" "Everything up-to-date" "$0" $LINENO <<'_eof'
1..7
ok 1 - No files with unexpected DEBUG comments
ok 2 - No files with unexpected TODO comments
ok 3 - No TAP_TODO's
ok 4 - Not using github.com/magisterquis/mqd
ok 5 - Running with -h exits happily
ok 6 - Packages up-to-date
ok 7 - Latest Go version will be used
_eof

# vim: ft=sh
