// Program toolskel generates boilerplate for small tools written in Go.
package main

/*
 * toolskel.go
 * Generate command boilerplate
 * By J. Stuart McMurray
 * Created 20230204
 * Last Modified 20250201
 */

import (
	"cmp"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/magisterquis/toolskel/internal/gencode"
	"golang.org/x/tools/txtar"
)

// Default filenames
var (
	GitignoreName   = ".gitignore"
	MakefileName    = "Makefile"
	ReadmeName      = "README.md"
	StaticcheckName = "staticcheck.conf"
)

func main() {
	var (
		createProgram = flag.Bool(
			"program",
			false,
			"Generate a main() program skeleton in name.go",
		)
		createLibrary = flag.Bool(
			"library",
			false,
			"Generate a library skeleton in name.go",
		)
		createMakefile = flag.Bool(
			"makefile",
			false,
			"Generate a Makefile",
		)
		createStaticcheck = flag.Bool(
			"staticcheck",
			false,
			"Generate a sensible staticcheck.conf",
		)
		createProgramReadme = flag.Bool(
			"program-readme",
			false,
			"Generate a README.md suitable for a program",
		)
		createLibraryReadme = flag.Bool(
			"library-readme",
			false,
			"Generate a README.md suitable for a library",
		)
		createGitignore = flag.Bool(
			"gitignore",
			false,
			"Generate a .gitignore",
		)
		/* Bulk creation. */
		newProgram = flag.Bool(
			"new-program",
			false,
			"Same as -program -gitignore -makefile "+
				"-program-readme -staticcheck",
		)
		newLibrary = flag.Bool(
			"new-library",
			false,
			"Same as -library -library-readme -staticcheck",
		)
		/* Other options. */
		dir = flag.String(
			"dir",
			".",
			"Directory in which to create files",
		)
		overwrite = flag.Bool(
			"overwrite",
			false,
			"Overwrite existing files",
		)
		author = flag.String(
			"author",
			defaultUsername(),
			"Author's `name`",
		)
		toTxtar = flag.Bool(
			"txtar",
			false,
			"Write a txtar achive to stdout intead of files",
		)
		today = flag.String(
			"today",
			time.Now().Format("20060102"),
			"Created `date` for generated files",
		)
		quiet = flag.Bool(
			"quiet",
			false,
			"Only log errors",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %s [options] [name [description...]]

Generates boilerplate Go projects.  Go source files will be named name.go.

Options:
`,
			filepath.Base(os.Args[0]),
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* Be in the file-creation directory. */
	if !*toTxtar {
		if err := os.Chdir(*dir); nil != err {
			log.Fatalf(
				"Error setting working directory to %s: %s",
				*dir,
				err,
			)
		}
	}

	/* Set bulk options. */
	if *newProgram {
		*createGitignore = true
		*createMakefile = true
		*createProgram = true
		*createProgramReadme = true
		*createStaticcheck = true
	}
	if *newLibrary {
		*createLibrary = true
		*createLibraryReadme = true
		*createStaticcheck = true
	}

	/* Can create a program xor a library. */
	if *createProgram && *createLibrary {
		log.Fatalf("Cannot create both a program and a library")
	} else if *createProgramReadme && *createLibraryReadme {
		log.Fatalf("Cannot create both a program and a library README")
	}

	/* Generation parameters. */
	params := gencode.Params{
		Name:   flag.Arg(0),
		Author: *author,
		Today:  *today,
	}
	if 2 <= flag.NArg() {
		params.Description = strings.Join(flag.Args()[1:], " ")
	}

	/* Prep a txtar archive, if we're doing that. */
	var ta txtar.Archive
	if *toTxtar {
		ta.Comment = []byte(fmt.Sprintf(
			"Created by toolskel %s",
			time.Now().Format(time.RFC3339),
		))
	}

	/* Generate ALL the things. */
	type fileConf struct {
		do   bool
		name string
		gen  gencode.Generator
	}
	var (
		tool  = cmp.Or(params.Name, gencode.DefaultName)
		confs = []fileConf{{
			do:   *createProgram,
			name: tool + ".go",
			gen:  params.Program,
		}, {
			do:   *createLibrary,
			name: tool + ".go",
			gen:  params.Library,
		}, {
			do:   *createMakefile,
			name: MakefileName,
			gen:  params.Makefile,
		}, {
			do:   *createStaticcheck,
			name: StaticcheckName,
			gen:  params.Staticcheck,
		}, {
			do:   *createProgramReadme,
			name: ReadmeName,
			gen:  params.Programreadme,
		}, {
			do:   *createLibraryReadme,
			name: ReadmeName,
			gen:  params.Libraryreadme,
		}, {
			do:   *createGitignore,
			name: GitignoreName,
			gen:  params.Gitignore,
		}}
	)
	/* Make sure we actually have something to do. */
	if !slices.ContainsFunc(confs, func(fc fileConf) bool {
		return fc.do
	}) {
		log.Fatalf("Need something to generate")
	}
	/* Generate ALL the files. */
	for _, fc := range confs {
		/* Easy if we're not generating this one. */
		if !fc.do {
			continue
		}
		/* Writing to a file is easy. */
		if !*toTxtar {
			if err := params.ToFile(
				fc.name,
				fc.gen,
				*overwrite,
			); nil != err {
				log.Fatalf(
					"Error creating %s: %s",
					fc.name,
					err,
				)
			}
			if !*quiet {
				log.Printf("Created %s", fc.name)
			}
			continue
		}
		/* Update the archive. */
		var (
			tf  = txtar.File{Name: fc.name}
			err error
		)
		if tf.Data, err = fc.gen(); nil != err {
			log.Fatalf("Error generating %s: %s", fc.name, err)
		}
		ta.Files = append(ta.Files, tf)
		if !*quiet {
			log.Printf("Generated %s", fc.name)
		}
	}

	/* Write the archive, if we're archiving. */
	if *toTxtar {
		if _, err := os.Stdout.Write(txtar.Format(&ta)); nil != err {
			log.Fatalf("Error writing archive: %s", err)
		}
	}
}

// defaultUsername returns the current user's name or username, if available.
func defaultUsername() string {
	u, err := user.Current()
	if nil != err {
		log.Printf("Unable to get current user's info: %s", err)
		u = &user.User{}
	}
	for _, n := range []string{u.Name, u.Username} {
		if "" != n {
			return n
		}
	}
	return ""
}
