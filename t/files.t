#!/bin/ksh
#
# files.t
# Test file generation
# By J. Stuart McMurray
# Created 20250201
# Last Modified 20250407

set -uo pipefail

. ./t/shmore.subr

tap_plan 12

TESTDATA=t/testdata/files

# Update all the files with macros
make -C "$TESTDATA" -f ../../updatemacros.mk -s
tap_ok $? "Macro files up-to-date" "$0" $LINENO

# test_gen tests generating using the flag -$1 and checks agains the contents
# of the directory $2.
test_gen() {
        tap_plan 3

        # Temporary directory for files
        TD=$(mktemp -d)
        tap_ok "$?" "Made temporary directory" "$0" $LINENO
        trap 'rm -rf "$TD"; tap_done_testing' EXIT # Make sure it's removed

        # Generate the files
        go run . \
                "-$1" \
                -author "test_author" \
                -dir "$TD" \
                -name "" \
                -quiet \
                -today 25251225
        tap_ok $? "Ran successfully" "$0" $LINENO

        # See where they differ
        GOT=$(diff -ruN -x '*.m4' -x '*.swp' "$TD" "$DIR" ||:)
        tap_is "$GOT" "" "No differences found" "$0" $LINENO
}

# Test ALL the things
for DIR in t/testdata/files/*; do
        FLAG=${DIR##*/}
        subtest() { test_gen "$FLAG" "$DIR"; }
        tap_subtest "Generate -$FLAG" "subtest" "$0" $LINENO
done

# vim: ft=sh
