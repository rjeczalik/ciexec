package ciexec

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/rjeczalik/ciexec/bash"
)

var (
	all []interface {
		Is([]byte) bool
		Exec(string, io.Reader, io.Writer) error
	}
	defaultEnv []string = os.Environ()
)

// SplitEnv TODO
func SplitEnv(env []string) (m map[string]string) {
	m = make(map[string]string, len(env))
	for _, env := range env {
		if i := strings.Index(env, "="); i > 0 && i < len(env)-1 {
			m[env[:i]] = env[i+1:]
		}
	}
	return
}

// MergeEnv TODO
func MergeEnv(env, ovr []string) (mrg []string) {
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

// ExecCmd TODO
func ExecCmd(cmds []*exec.Cmd, w io.Writer) error {
	var err error
	for _, cmd := range cmds {
		cmd.Stdout, cmd.Stderr = w, w
		if e := cmd.Run(); e != nil && err == nil {
			err = e
		}
	}
	return err
}

// ExecBash TODO
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

// Exec TODO
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
		if ci.Is(b) {
			return ci.Exec(detail, bytes.NewBuffer(b), w)
		}
	}
	return nil
}
