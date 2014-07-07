// cmd/ciexec executes CI configuration file like it was a shell script.
//
// Current version of cmd/ciexec supports Travis' and Pulse's configuration
// files. Travis files (typically .travis.yml in the root of the repository)
// are executed as documented in:
//
//   http://docs.travis-ci.com/user/build-configuration/#Build-Lifecycle
//
// The build matrix settings are ignored, user is reposinbile for ensuring
// the environment is capable of executing every commands. In order to execute
// a .travis.yml file run:
//
//   ~ $ ciexec travis .travis.yml
//
// cmd/ciexec needs a recipe name to execute a Pulse file, e.g.:
//
//   ~ $ ciexec go .pulse.xml
//
// Usage
//
//   NAME:
//     ciexec - executes CI configuration file like it was a shell script
//
//   USAGE:
//     ciexec <pulse recipe> [<versioned pulse file>]
//     ciexec travis [<travis configuration file>]
//
//   EXAMPLES:
//     ~ $ ciexec travis src/github.com/rjeczalik/ciexec/.travis.yml
//     ~ $ ciexec go src/github.com/rjeczalik/ciexec/.pulse.xml
package main

import (
	"fmt"
	"os"

	"github.com/rjeczalik/ciexec"
)

const usage = `NAME:
	ciexec - executes CI configuration file like it was a shell script

USAGE:
	ciexec <pulse recipe> [<versioned pulse file>]
	ciexec travis [<travis configuration file>]

EXAMPLES:
	~ $ ciexec travis src/github.com/rjeczalik/ciexec/.travis.yml
	~ $ ciexec go src/github.com/rjeczalik/ciexec/.pulse.xml`

func die(v interface{}) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}

func ishelp(s string) bool {
	return s == "-h" || s == "-help" || s == "help" || s == "--help" || s == "/?"
}

func main() {
	if len(os.Args) > 3 {
		die(usage)
	}
	if len(os.Args) > 1 && ishelp(os.Args[1]) {
		fmt.Println(usage)
		return
	}
	var file, detail string
	if len(os.Args) == 3 {
		file = os.Args[2]
	}
	if len(os.Args) > 1 {
		detail = os.Args[1]
	}
	if err := ciexec.Exec(file, detail, os.Stdout); err != nil {
		die(err)
	}
}
