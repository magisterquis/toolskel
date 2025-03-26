#!/bin/ksh
#
# makefile.t
# Make sure our Makefile is the same as what's generated
# By J. Stuart McMurray
# Created 20250222
# Last Modified 20250326

set -uo pipefail

. ./t/shmore.subr

tap_plan 8

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

# Make sure the long line check works.
cp t/testdata/makefile/long_usage_line.go "$TMPD"
GOT=$(cd "$TMPD" && go mod init makefile_t 2>/dev/null && make -i gotest)
WANT="$(cat <<'_eof'
go test -trimpath -ldflags "-w -s" -timeout 3s ./...
?   	makefile_t	[no test files]
go vet -trimpath -ldflags "-w -s" ./...
staticcheck ./...
go run -trimpath -ldflags "-w -s" . -h 2>&1 | awk ' /^Options:$|MQD DEBUG PACKAGE LOADED$/ { exit } /^Usage: / { sub(/^Usage: [^[:space:]]+\//, "Usage: ") } /.{80,}/ { print "Long usage line: " $0; exit 1 } '
Long usage line: In taberna quando sumus non curamus quid sit humus sed ad ludum properamus cui semper insudamus.
*** Error 1 in target 'gotest' (ignored)
_eof
)"
tap_is "$GOT" "$WANT" "Long usage line detected" "$0" $LINENO

# vim: ft=sh
