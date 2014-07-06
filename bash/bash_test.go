package bash

import (
	"bytes"
	"strings"
	"testing"
)

func TestBash(t *testing.T) {
	t.Parallel()
	cases := [...]struct {
		cmd []string
		exp []string
		err bool
	}{
		0: {
			[]string{"echo DU", "echo PA", "echo 123"},
			[]string{"DU", "PA", "123"},
			false,
		},
		1: {
			[]string{"for i in {1..10}; do echo -n $i; done"},
			[]string{"12345678910"},
			false,
		},
		2: {
			[]string{"TEST_VAR=XD", "for i in {1..3}; do echo $TEST_VAR; done"},
			[]string{"XD", "XD", "XD"},
			false,
		},
		3: {
			[]string{"echo DU", "ercdfg PA", "echo 123"},
			[]string{"DU", "123"},
			true,
		},
		4: {
			[]string{"for i in {1..10}; do echo -n $i; done", "exit 10"},
			[]string{"12345678910"},
			true,
		},
		5: {
			[]string{"TEST_VAR=echro", "eval $TEST_VAR"},
			nil,
			true,
		},
		6: {
			[]string{"TEST_CMD='echo ABC'", "TEST_VAR=$($TEST_CMD)", "echo DGF", "false"},
			[]string{"ABC", "DGF"},
			true,
		}}
	for i, cas := range cases {
		var b bytes.Buffer
		s, err := NewSession(&b)
		if err != nil {
			t.Errorf("want err=nil; got %q (i=%d)", err, i)
			continue
		}
		for _, cmd := range cas.cmd {
			s.Start(cmd)
		}
		if err = s.Close(); !cas.err && err != nil {
			t.Errorf("want err=nil; got %q (i=%d)", err, i)
			continue
		}
		if cas.err && err == nil {
			t.Errorf("want err!=nil (i=%d)", i)
		}
		out, j := b.String(), 0
		for _, exp := range cas.exp {
			if k := strings.Index(out[j:], exp); k == -1 {
				t.Errorf("missing %q in the output (i=%d)", exp, i)
			} else {
				j = k + len(exp)
			}
		}
	}
}
