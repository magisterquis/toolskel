toolskel
========
Generates boilerplate for small tools written in Go.

Expected userbase size: `1`

Example output can be found in
[`t/testdata/files`](./t/testdata/files).

Quickstart
----------
1. Have a working Go installation.
2. Grab the code.
   ```sh
   go install github.com/magisterquis/toolskel@latest
   ```
3. Be in a directory where a new project will live.
4. Create some files.
   ```sh
   toolskel -new-program thing A cool Thing
   ```
5. Code away!
   ```sh
   vim -p README.md thing.go
   make test
   ```

Usage
-----
```
Usage: toolskel [options] [description...]

Generates boilerplate Go projects.  Go source files will be named name.go.

Options:
  -author name
    	Author's name (default "J. Stuart McMurray")
  -basic-tests
    	Generate a ./t directory with basic tests
  -dir string
    	Directory in which to create files (default ".")
  -gitignore
    	Generate a .gitignore
  -library
    	Generate a library skeleton in name.go
  -library-readme
    	Generate a README.md suitable for a library
  -makefile
    	Generate a Makefile
  -name name
    	Project name (default "toolskel")
  -new-library
    	Same as -library -library-readme -staticcheck
  -new-program
    	Same as -basic-tests -program -gitignore -makefile -program-readme -shmore -staticcheck
  -overwrite
    	Overwrite existing files
  -program
    	Generate a main() program skeleton in name.go
  -program-readme
    	Generate a README.md suitable for a program
  -quiet
    	Only log errors
  -shmore
    	Generate a ./t directory with shmore
  -staticcheck
    	Generate a sensible staticcheck.conf
  -today date
    	Created date for generated files (default "20250310")
  -txtar
    	Write a txtar achive to stdout intead of files
```
