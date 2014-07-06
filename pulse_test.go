package ciexec

import (
	"bytes"
	"encoding/xml"
	"reflect"
	"testing"
)

func cmdequal(lhs, rhs Commands) bool {
	if len(lhs) != len(rhs) {
		return false
	}
	for k := range rhs {
		if _, ok := lhs[k]; !ok {
			return false
		}
		if len(rhs[k]) != len(lhs[k]) {
			return false
		}
		for i := range rhs[k] {
			if rhs[k][i].Path != lhs[k][i].Path {
				return false
			}
			if len(rhs[k][i].Args) != len(lhs[k][i].Args) {
				return false
			}
			if !reflect.DeepEqual(rhs[k][i].Args, lhs[k][i].Args) {
				return false
			}
		}
	}
	return true
}

var fixture = [...]struct {
	cfg []byte
	pro *Project
	cmd Commands
}{
	0: {
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<project>
	<recipe name="test">
		<executable args="-arg1 -arg2 -arg3" exe="jkll"/>
		<executable args="-merg1 -merg2 -merg3" exe="meho">
			<arg>long merg 4</arg>
			<arg>long merg 5</arg>
		</executable>
	</recipe>
</project>`),
		&Project{
			XMLName: xml.Name{Local: "project"},
			Recipe: []Recipe{{
				Name: "test",
				Executable: []Executable{{
					Args: "-arg1 -arg2 -arg3",
					Exe:  "jkll",
				}, {
					Args:    "-merg1 -merg2 -merg3",
					Exe:     "meho",
					ArgList: []string{"long merg 4", "long merg 5"},
				}},
			}},
		},
		Commands{
			"test": {
				{Path: "jkll", Args: []string{"jkll", "-arg1", "-arg2", "-arg3"}},
				{Path: "meho", Args: []string{"meho", "-merg1", "-merg2", "-merg3", "long merg 4", "long merg 5"}},
			},
		},
	},
	1: {
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<project>
	<recipe name="test1">
		<executable args="-arg1 -arg2" exe="/usr/bin/zywww"/>
		<executable args="-merg1 -merg2" exe="/sbin/meho">
			<arg>- m e r g 3</arg>
			<arg>- m e r g 4</arg>
		</executable>
	</recipe>
	<recipe name="test2">
		<executable args="-a1 -a2" exe="/bin/xyz"/>
		<executable args="-m1 -m2" exe="/usr/local/bin/meho">
			<arg>- m 3</arg>
			<arg>- m 4</arg>
		</executable>
		<executable args="" exe="z">
			<arg>- z 1</arg>
			<arg></arg>
			<arg>- z 2</arg>
		</executable>
		<executable args="" exe="x">
			<arg></arg>
			<arg></arg>
		</executable>
	</recipe>
</project>`),
		&Project{
			XMLName: xml.Name{Local: "project"},
			Recipe: []Recipe{{
				Name: "test1",
				Executable: []Executable{{
					Args: "-arg1 -arg2",
					Exe:  "/usr/bin/zywww",
				}, {
					Args:    "-merg1 -merg2",
					Exe:     "/sbin/meho",
					ArgList: []string{"- m e r g 3", "- m e r g 4"},
				}},
			}, {
				Name: "test2",
				Executable: []Executable{{
					Args: "-a1 -a2",
					Exe:  "/bin/xyz",
				}, {
					Args:    "-m1 -m2",
					Exe:     "/usr/local/bin/meho",
					ArgList: []string{"- m 3", "- m 4"},
				}, {
					Exe:     "z",
					ArgList: []string{"- z 1", "- z 2"},
				}, {
					Exe: "x",
				}},
			}},
		},
		Commands{
			"test1": {
				{Path: "/usr/bin/zywww", Args: []string{"/usr/bin/zywww", "-arg1", "-arg2"}},
				{Path: "/sbin/meho", Args: []string{"/sbin/meho", "-merg1", "-merg2", "- m e r g 3", "- m e r g 4"}},
			},
			"test2": {
				{Path: "/bin/xyz", Args: []string{"/bin/xyz", "-a1", "-a2"}},
				{Path: "/usr/local/bin/meho", Args: []string{"/usr/local/bin/meho", "-m1", "-m2", "- m 3", "- m 4"}},
				{Path: "z", Args: []string{"z", "- z 1", "- z 2"}},
				{Path: "x", Args: []string{"x"}},
			},
		},
	},
}

func TestPulseIs(t *testing.T) {
	t.Parallel()
	for i, cas := range fixture {
		if !(pulse{}).Is(cas.cfg) {
			t.Errorf("(pulse{}).Is failed (i=%d)", i)
		}
	}
}

func TestPulseParse(t *testing.T) {
	t.Parallel()
	for i, cas := range fixture {
		pro, err := pulseParse(bytes.NewReader(cas.cfg))
		if err != nil {
			t.Errorf("want err=nil; got %q (i=%d)", err, i)
			continue
		}
		if !reflect.DeepEqual(pro, cas.pro) {
			t.Errorf("want pro=%+v; got %+v (i=%d)", cas.pro, pro, i)
		}
	}
}

func TestPulseCmd(t *testing.T) {
	t.Parallel()
	for i, cas := range fixture {
		cmd := pulseCmd(cas.pro)
		if !cmdequal(cmd, cas.cmd) {
			t.Errorf("want cmd=%+v; got %+v (i=%d)", cas.cmd, cmd, i)
		}
	}
}
