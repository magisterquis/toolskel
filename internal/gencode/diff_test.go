package gencode

/*
 * diff_test.go
 * Testing whether two files are the same
 * By J. Stuart McMurray
 * Created 20230415
 * Last Modified 20230425
 */

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
	"text/tabwriter"
)

// nContextLines is the number of lines before and after a mismatch to print.
const nContextLines = 3

// lineSep separates lines, as a []byte.
var lineSep = []byte("\n")

// errorIfDiff returns true if got and want don't have the same bytes.  Before
// it returns true, it calls t.Errorf to report the error.  If the slices are
// the same, t.Error is not called.  gotN and wantN are used in place of "got"
// and "want" in the error message, if one is to be printed.  If either is the
// empty string, a sensible default ("got" or "want") will be used.
func errorIfDiff(t *testing.T, got, want []byte, gotN, wantN string) bool {
	t.Helper()

	/* Figure out where they differ. */
	ln := diff(t, got, want)
	if 0 == ln {
		return false
	}

	/* Make sure we have got/want names. */
	if "" == gotN {
		gotN = "got"
	}
	if "" == wantN {
		wantN = "want"
	}

	/* Roll a nice message. */
	t.Errorf(
		"discrepancy found at line %d:\n"+
			"--- %s (+/- %d):\n%s"+
			"--- %s (+/- %d):\n%s",
		ln,
		gotN, nContextLines, contextLines(got, ln),
		wantN, nContextLines, contextLines(want, ln),
	)
	return true
}

func TestDiff(t *testing.T) {
	for _, c := range []struct {
		a    string
		b    string
		want int
	}{{
		a:    "a\nb\nc\nd\ne\nf\ng\nh\ni\nj\n",
		b:    "a\nb\nc\nd\ne\nf\ng\nh\ni\nj\n",
		want: 0,
	}, {
		a:    "a\nb\nc\nd\ne\nf\ng\nh\ni\nj\n",
		b:    "b\nb\nc\nd\ne\nf\ng\nh\ni\nj\n",
		want: 1,
	}, {
		a:    "a\nb\nc\nd\ne\nf\ng\nh\ni\nj\n",
		b:    "a\nb\nc\nd\ne\nf\ng\nh\ni\nx\n",
		want: 10,
	}} {
		c := c /* :S */

		t.Run("", func(t *testing.T) {
			t.Parallel()
			got := diff(t, []byte(c.a), []byte(c.b))
			if got != c.want {
				t.Errorf("got: %d", got)
				return
			}
		})
	}

}

// diff assumes a and b are text, and returns the first 1-indexed line number
// where they differ.  If the two slices are the same, diff return 0.
func diff(t *testing.T, a, b []byte) int {
	t.Helper()

	/* If the slices are the same, life's easy. */
	if bytes.Equal(a, b) {
		return 0
	}

	/* If exactly one's empty, it's on the first line. */
	if 0 == len(a) || 0 == len(b) {
		return 1
	}

	/* Get the bytes as lines. */
	als := bytes.Split(a, lineSep)
	bls := bytes.Split(b, lineSep)

	/* Compare line-by-line until we find the discrepancy. */
	for i, al := range als {
		/* If we haven't any more lines in b, the discrepany is
		here. */
		if len(bls) <= i {
			return len(bls) - 1
		}

		if !bytes.Equal(al, bls[i]) {
			return i + 1
		}
	}

	/* If we're here, we ran out of lines in a before we finished b. */
	return len(als) - 1
}

func TestContextLines(t *testing.T) {
	for _, c := range []struct {
		Have string
		N    int
		Want string
	}{{
		Have: "a\nb\nc\nd\ne\nf\ng\nh\ni\nj\n",
		N:    5,
		Want: "2  b\n3  c\n4  d\n5  e <--\n6  f\n7  g\n8  h\n",
	}, {
		Have: "a\nb\nc\nd\ne\nf\ng\nh\ni\nj\n",
		N:    1,
		Want: "1  a <--\n2  b\n3  c\n4  d\n",
	}, {
		Have: "a\nb\nc\nd\ne\nf\ng\nh\ni\nj\n",
		N:    2,
		Want: "1  a\n2  b <--\n3  c\n4  d\n5  e\n",
	}, {
		Have: "a\nb\nc\nd\ne\nf\ng\nh\ni\nj\n",
		N:    9,
		Want: "6   f\n7   g\n8   h\n9   i <--\n10  j\n",
	}, {
		Have: "a\nb\nc\nd\ne\nf\ng\nh\ni\nj\n",
		N:    10,
		Want: "7   g\n8   h\n9   i\n10  j <--\n",
	}} {
		c := c /* :| */
		t.Run(strconv.Itoa(c.N), func(t *testing.T) {
			got := contextLines([]byte(c.Have), c.N)
			if string(got) == c.Want {
				return
			}
			t.Errorf("\nWant:%q,", got)
		})
	}
}

// contextLines returns line n in b, plus nContextLines lines before and after
// and line numbers.  The mismatched line is noted with an <--.
func contextLines(b []byte, n int) []byte {
	if 0 >= n {
		panic(fmt.Sprintf("contextLines would look for %d", n))
	}
	n-- /* Input's 1-indexed. */

	/* Split into lines. */
	ls := bytes.Split(bytes.TrimSuffix(b, lineSep), lineSep)

	/* Work out the first and last line to get. */
	start := n - nContextLines
	if 0 > start {
		start = 0
	}
	end := n + nContextLines + 1
	if len(ls) < end {
		end = len(ls)
	}
	badOff := n - start

	/* Pretty up the lines to return. */
	var (
		buf bytes.Buffer
		tw  = tabwriter.NewWriter(&buf, 2, 8, 2, ' ', 0)
	)
	for i, l := range ls[start:end] {
		if i == badOff {
			l = append(l, []byte(" <--")...)
		}
		fmt.Fprintf(tw, "%d\t%s\n", i+start+1, l)
	}
	tw.Flush()

	return buf.Bytes()
}
