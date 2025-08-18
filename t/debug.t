#!/bin/ksh
#
# debug.t
# Make sure we catch debugging statements
# By J. Stuart McMurray
# Created 20250818
# Last Modified 20250818

set -euo pipefail

. ./t/shmore.subr

tap_plan 3

TMPD=$(mktemp -d)
mkdir -p "$TMPD/t"
TD=$(basename "$0")
TD=t/testdata/${TD%.t}
trap 'rm -rf "$TMPD"; tap_done_testing' exit

# Copy our test file.
set -A IFNS $(find "$TD" -name '*.b64' -type f) 
for FN in "${IFNS[@]}"; do 
        OFN=${FN#$TD/}
        OFN="$TMPD/${OFN%.b64}"
        mkdir -p "$(dirname "$OFN")"
        openssl base64 -A -d <"$FN" >"$OFN"
done
WANT=${#IFNS[@]}
GOT=$(($(find "$TMPD" -type f | wc -l)))
tap_is "$GOT" "$WANT" "Decoded all of the test files" "$0" $LINENO

# Set up to run basic tests
go run .  -basic-tests -dir "$TMPD" -quiet
tap_ok "$?" "Created files happily" "$0" $LINENO
cp t/shmore.subr "$TMPD"/t

# See if it worked.
GOT=$(cd "$TMPD" && tap_reset && ksh t/basic_tests.t 2>&1 ||:)
WANT="$(cat <<'_eof' | sed -e 's/DExBUG/DEBUG/g' -e 's/TOxDO/TODO/g'
1..7
not ok 1 - No files with unexpected DExBUG comments
#   Failed test 'No files with unexpected DExBUG comments'
#   at t/basic_tests.t line 28.
#          got: 'debug_c:2:/* DExBUG: Catch this. */
# debug_c:3:Another line /* DExBUG */
# debug_shell:2:# DExBUG: Catch this
# debug_shell:3:Another line # DExBUG: Catch this
# debug_text:2:DExBUG: Catch this'
#     expected: ''
not ok 2 - No files with unexpected TOxDO comments
#   Failed test 'No files with unexpected TOxDO comments'
#   at t/basic_tests.t line 35.
#          got: 'todo_c:2:/* TOxDO: Catch this. */
# todo_shell:2:# TOxDO: Catch this
# todo_text:2:TOxDO: Catch this'
#     expected: ''
not ok 3 - No TAP_TOxDO's
#   Failed test 'No TAP_TOxDO's'
#   at t/basic_tests.t line 37.
#          got: 't/tap_todo.t:7:TAP_TOxDO="A TAP TOxDO"
# t/tap_todo.t:8:        TAP_TOxDO="An indented TAP TOxDO"'
#     expected: ''
ok 4 # skip Not a Go program
ok 5 # skip Not a Go program
ok 6 # skip Not a Go program
ok 7 # skip Not a Go program
# Looks like you failed 3 tests of 7.
_eof
)"

# Add a formatting line if we're running in a harness.
if [[ -n "${HARNESS_ACTIVE-}" ]]; then
        WANT=$(print -r "$WANT" | awk '1;/^not ok [[:digit:]]/{print""}')
fi

print -r "$GOT" >$TMPD/got
print -r "$WANT" >$TMPD/want
DIFF=$(diff -u "$TMPD/want" "$TMPD/got" 2>&1 ||:)
tap_is "$DIFF" "" "Found the correct DEBUG/TODO lines" "$0" $LINENO

# vim: ft=sh
