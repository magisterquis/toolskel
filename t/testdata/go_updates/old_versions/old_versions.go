package main

/*
 * old_versions.go
 * Dummy program to load libraries
 * By J. Stuart McMurray
 * Created 20250407
 * Last Modified 20250407
 */

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/net/xsrftoken"
	"golang.org/x/text/cases"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "A dummy usage message\n")
	}
	flag.Parse()
	fmt.Printf("XSRF token timeout: %s", xsrftoken.Timeout)
	fmt.Printf("Unicode version: %s", cases.UnicodeVersion)
}
