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

func (pulse) Is(content []byte) bool {
	return strings.Contains(http.DetectContentType(content), "text/xml")
}

func (pulse) Exec(recipe string, r io.Reader, w io.Writer) error {
	pro, err := pulseParse(r)
	if err != nil {
		return err
	}
	rec := pulseCmd(pro)
	if rec == nil {
		return errors.New("no recipes found")
	}
	cmd, ok := rec[recipe]
	if !ok {
		return errors.New("no recipe found or it has no commands: " + recipe)
	}
	return ExecCmd(cmd, w)
}

// Commands TODO
type Commands map[string][]*exec.Cmd

// Executable TODO
type Executable struct {
	Args    string   `xml:"args,attr"`
	Exe     string   `xml:"exe,attr"`
	ArgList []string `xml:"arg"`
}

// Recipw TODO
type Recipe struct {
	Name       string       `xml:"name,attr"`
	Executable []Executable `xml:"executable"`
}

// Project TODO
type Project struct {
	XMLName xml.Name `xml:"project"`
	Recipe  []Recipe `xml:"recipe"`
}

func pulseParse(r io.Reader) (*Project, error) {
	var (
		dec = xml.NewDecoder(r)
		pro = &Project{}
	)
	if err := dec.Decode(pro); err != nil {
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

func pulseCmd(pro *Project) Commands {
	var m = make(Commands)
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
