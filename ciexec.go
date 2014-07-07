// Package ciexec implements functions to execute some CI configuration
// scripts.
package ciexec

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/rjeczalik/ciexec/bash"
)

var (
	all []interface {
		Is(string, []byte) bool
		Exec(string, []byte, io.Writer) error
	}
	defaultEnv []string = os.Environ()
)

// SplitEnv splits key=value pairs from env into map, filtering out env with
// empty values.
func SplitEnv(env []string) (m map[string]string) {
	m = make(map[string]string, len(env))
	for _, env := range env {
		if i := strings.Index(env, "="); i > 0 && i < len(env)-1 {
			m[env[:i]] = env[i+1:]
		}
	}
	return
}

// OverrideEnv overrides environment variables stored in env with ones from ovr.
// It ignores environment variables with empty values.
func OverrideEnv(env, ovr []string) (mrg []string) {
	tmp := SplitEnv(env)
	for k, v := range SplitEnv(ovr) {
		tmp[k] = v
	}
	mrg = make([]string, 0, len(tmp))
	for k, v := range tmp {
		mrg = append(mrg, k+"="+v)
	}
	return
}

// ExecCmd runs commands in order they're stored in cmds, writing combined output
// from each command to w. Execution loop does not stop even if a command fails.
// The function returns the last error encountered, if any.
func ExecCmd(cmds []*exec.Cmd, w io.Writer) error {
	var err error
	for _, cmd := range cmds {
		cmd.Stdout, cmd.Stderr = w, w
		if e := cmd.Run(); e != nil {
			err = e
		}
	}
	return err
}

// ExecBash runs script expressions in a single bash session. Execution loop
// does not stop, even if one of the commands exits with a code != 0.
// The function returns the last error encountered, if any.
func ExecBash(cmds []string, w io.Writer) error {
	s, err := bash.NewSession(w)
	if err != nil {
		return err
	}
	for _, cmd := range cmds {
		s.Start(cmd)
	}
	return s.Close()
}

// Exec detects the file's CI format type and executes it passing CI-specific
// detail string. It writes outputs from any commands it executes to the w.
func Exec(file, detail string, w io.Writer) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	for _, ci := range all {
		if ci.Is(file, b) {
			return ci.Exec(detail, b, w)
		}
	}
	return nil
}
