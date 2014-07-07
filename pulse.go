package ciexec

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

func init() {
	all = append(all, pulse{})
}

type pulse struct{}

func (pulse) Is(_ string, content []byte) bool {
	return strings.Contains(http.DetectContentType(content), "text/xml")
}

func (pulse) Exec(recipe string, p []byte, w io.Writer) error {
	pro, err := pulseParse(p)
	if err != nil {
		return err
	}
	rec := pulseCmd(pro)
	if rec == nil {
		return errors.New("pulse: no recipes found")
	}
	cmd, ok := rec[recipe]
	if !ok {
		return errors.New("pulse: no recipe found or it has no commands: " + recipe)
	}
	return ExecCmd(cmd, w)
}

type pulseCommands map[string][]*exec.Cmd

type (
	pulseProject struct {
		XMLName xml.Name `xml:"project"`
		Recipe  []Recipe `xml:"recipe"`
	}
	Recipe struct {
		Name       string       `xml:"name,attr"`
		Executable []Executable `xml:"executable"`
	}
	Executable struct {
		Args    string   `xml:"args,attr"`
		Exe     string   `xml:"exe,attr"`
		ArgList []string `xml:"arg"`
	}
)

func pulseParse(p []byte) (*pulseProject, error) {
	var pro = &pulseProject{}
	if err := xml.Unmarshal(p, pro); err != nil {
		return nil, err
	}
	// Filter out empty empty ArgList.
	for i := range pro.Recipe {
		for j := range pro.Recipe[i].Executable {
			tmp := &pro.Recipe[i].Executable[j].ArgList
			for k := 0; k < len(*tmp); {
				if (*tmp)[k] == "" {
					if len(*tmp) != k+1 {
						*tmp = append((*tmp)[:k], (*tmp)[k+1:]...)
					} else {
						*tmp = (*tmp)[:k]
					}
				} else {
					k += 1
				}
			}
			if len(*tmp) == 0 {
				*tmp = nil
			}
		}
	}
	return pro, nil
}

func pulseCmd(pro *pulseProject) pulseCommands {
	var m = make(pulseCommands)
	for i := range pro.Recipe {
		var c = make([]*exec.Cmd, 0, len(pro.Recipe[i].Executable))
		for j := range pro.Recipe[i].Executable {
			cmd := exec.Command(pro.Recipe[i].Executable[j].Exe)
			// TODO(rjeczalik): For arguments that needs quoting project config uses
			//                  <arg> child nodes, but it may be possible to pass quoted ones
			//                  to the args attribute.
			arg := strings.Split(pro.Recipe[i].Executable[j].Args, " ")
			for _, arg := range append(arg, pro.Recipe[i].Executable[j].ArgList...) {
				if arg != "" {
					cmd.Args = append(cmd.Args, arg)
				}
			}
			c = append(c, cmd)
		}
		if len(c) != 0 {
			m[pro.Recipe[i].Name] = c
		}
	}
	if len(m) != 0 {
		return m
	}
	return nil
}
