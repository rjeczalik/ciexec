package ciexec

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func testCommand(cmds, envs [][]string) (c []*exec.Cmd) {
	c = make([]*exec.Cmd, 0, len(cmds))
	for i := range cmds {
		cmd := exec.Command(os.Args[0])
		cmd.Args = append([]string{"-test.run=TestHelperCommand", "--"}, cmds[i]...)
		cmd.Env = []string{"TEST_HELPER_COMMAND=1"}
		c = append(c, cmd)
		if len(envs) > i {
			c[len(c)-1].Env = append(c[len(c)-1].Env, envs[i]...)
		}
	}
	return
}

func TestHelperCommand(t *testing.T) {
	if os.Getenv("TEST_HELPER_COMMAND") != "" {
		args, n, cmd := os.Args, 0, ""
		for i := range args {
			if args[i] == "--" && i < len(args)-2 {
				n = i
				break
			}
		}
		if n == 0 {
			fmt.Fprintln(os.Stderr, "invalid arguments")
			os.Exit(1)
		}
		cmd, args = args[n+1], args[n+2:]
		switch cmd {
		case "echo":
			v := make([]interface{}, 0, len(args))
			for _, arg := range args {
				v = append(v, os.ExpandEnv(arg))
			}
			fmt.Println(v...)
			os.Exit(0)
		case "exit":
			n, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, "invalid arguments")
				os.Exit(2)
			}
			os.Exit(n)
		default:
			fmt.Fprintln(os.Stderr, "invalid arguments")
			os.Exit(3)
		}
	}
}
