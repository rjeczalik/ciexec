package ciexec

import (
	"bytes"
	"os/exec"
	"reflect"
	"testing"
)

func envequal(lhs, rhs []string) bool {
	if len(lhs) != len(rhs) {
		return false
	}
NEXT:
	for _, lhs := range lhs {
		for _, rhs := range rhs {
			if rhs == lhs {
				continue NEXT
			}
		}
		return false
	}
	return true
}

func TestOverrideEnv(t *testing.T) {
	t.Parallel()
	cases := [...]struct {
		env []string
		ovr []string
		mrg []string
	}{{
		[]string{"A=a", "B=b", "C=c", "D=d"},
		[]string{"E=e"},
		[]string{"A=a", "B=b", "C=c", "D=d", "E=e"},
	}, {
		[]string{"A=a", "B=b", "C=c", "D=d"},
		[]string{"E=e", "A=x", "F=f"},
		[]string{"A=x", "B=b", "C=c", "D=d", "E=e", "F=f"},
	}, {
		[]string{"GOPATH=/Users/john", "PATH=/usr/bin:/bin", "USER=john"},
		[]string{"GOPATH=/Users/mike", "USER=mike"},
		[]string{"GOPATH=/Users/mike", "PATH=/usr/bin:/bin", "USER=mike"},
	}, {
		[]string{},
		[]string{"A=a", "B=b"},
		[]string{"A=a", "B=b"},
	}, {
		[]string{"A", "B=b", "C=", "D=d"},
		[]string{"E=e"},
		[]string{"B=b", "D=d", "E=e"},
	}, {
		[]string{"A", "B=b", "C=c", "D="},
		[]string{"E=e", "F=f", "G=", "G"},
		[]string{"B=b", "C=c", "E=e", "F=f"},
	}}
	for i, cas := range cases {
		if mrg := OverrideEnv(cas.env, cas.ovr); !envequal(mrg, cas.mrg) {
			t.Errorf("want mrg=%v; got %v (i=%d)", cas.mrg, mrg, i)
		}
	}
}

func TestExecCmd(t *testing.T) {
	t.Parallel()
	cases := [...]struct {
		cmds [][]string
		envs [][]string
		out  []byte
		err  interface{}
	}{
		0: {
			cmds: [][]string{
				{"echo", "a", "b", "c"},
			},
			out: []byte("a b c\n"),
		},
		1: {
			cmds: [][]string{
				{"echo", "a"},
				{"echo", "b"},
				{"echo", "c"},
			},
			out: []byte("a\nb\nc\n"),
		},
		2: {
			cmds: [][]string{
				{"echo", "a", "d", "f"},
				{"echo", "d", "b", "e"},
				{"echo", "f", "e", "c"},
			},
			out: []byte("a d f\nd b e\nf e c\n"),
		},
		3: {
			cmds: [][]string{
				{"echo", "a", "d", "f"},
				{"exit", "0"},
				{"echo", "d", "b", "e"},
				{"exit", "0"},
				{"exit", "0"},
				{"echo", "f", "e", "c"},
			},
			out: []byte("a d f\nd b e\nf e c\n"),
		},
		4: {
			cmds: [][]string{
				{"echo", "a", "b", "c"},
				{"exit", "127"},
			},
			out: []byte("a b c\n"),
			err: (*exec.ExitError)(nil),
		},
		5: {
			cmds: [][]string{
				{"echo", "a"},
				{"echo", "b"},
				{"exit", "1"},
				{"echo", "c"},
			},
			out: []byte("a\nb\nc\n"),
			err: (*exec.ExitError)(nil),
		},
		6: {
			cmds: [][]string{
				{"exit", "17"},
				{"echo", "a", "d", "f"},
				{"echo", "d", "b", "e"},
				{"echo", "f", "e", "c"},
			},
			out: []byte("a d f\nd b e\nf e c\n"),
			err: (*exec.ExitError)(nil),
		},
		7: {
			cmds: [][]string{
				{"exit", "17"},
				{"echo", "a", "d", "f"},
				{"echo", "d", "b", "e"},
				{"exit", "18"},
				{"exit", "0"},
				{"exit", "19"},
				{"echo", "f", "e", "c"},
				{"exit", "20"},
			},
			out: []byte("a d f\nd b e\nf e c\n"),
			err: (*exec.ExitError)(nil),
		},
	}
	for i, cas := range cases {
		var b bytes.Buffer
		err := ExecCmd(testCommand(cas.cmds, cas.envs), &b)
		if cas.err == nil && err != nil {
			t.Errorf("want err=nil; got %q (i=%d, b=%q)", err, i, b.Bytes())
			continue
		}
		if cas.err != nil && err == nil {
			t.Errorf("want err!=nil; (i=%d, b=%q)", i, b.Bytes())
			continue
		}
		if cas.err != nil && reflect.TypeOf(cas.err) != reflect.TypeOf(err) {
			t.Errorf("want typeof(err)=%T; got %T (i=%d, b=%q)", cas.err, err, i, b.Bytes())
		}
		if !bytes.Equal(b.Bytes(), cas.out) {
			t.Errorf("want b=%q; got %q (i=%d)", cas.out, b.Bytes(), i)
		}
	}
}
