// TODO
package main

import (
	"fmt"
	"os"

	"github.com/rjeczalik/ciexec"
)

const usage = `NAME:
	ciexec - executes CI configuration file like it was a shell script

USAGE:
	ciexec recipe [versioned pulse file]`

func die(v interface{}) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}

func ishelp(s string) bool {
	return s == "-h" || s == "-help" || s == "help" || s == "--help" || s == "/?"
}

func main() {
	if len(os.Args) != 2 && len(os.Args) != 3 {
		die(usage)
	}
	if ishelp(os.Args[1]) {
		fmt.Println(usage)
		return
	}
	file := ".pulse.xml"
	if len(os.Args) == 3 {
		file = os.Args[2]
	}
	if err := ciexec.Exec(file, os.Args[1], os.Stdout); err != nil {
		die(err)
	}
}
