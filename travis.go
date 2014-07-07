package ciexec

import (
	"errors"
	"io"
	"path/filepath"
	"sort"

	"github.com/rjeczalik/ciexec/bash"
	"gopkg.in/v1/yaml"
)

func init() {
	all = append(all, travis{})
}

type travis struct{}

func (travis) Is(file string, _ []byte) bool {
	return filepath.Base(file) == ".travis.yml"
}

// TODO(rjeczalik): Support matrix builds, use 'detail' to configure local execution. (#2)
// TODO(rjeczalik): Support after_success/after_failure stages. (#1)
func (travis) Exec(detail string, p []byte, w io.Writer) error {
	pro, err := travisParse(p)
	if err != nil {
		return err
	}
	cmd := travisCmd(pro)
	if cmd == nil {
		return errors.New("travis: no commands found")
	}
	s, err := bash.NewSession(w)
	if err != nil {
		return err
	}
	for _, stage := range stages {
		if cmd, ok := cmd[stage]; ok {
			for _, cmd := range cmd {
				s.Start(cmd)
			}
		}
	}
	return s.Close()
}

var (
	// Keep this slice sorted.
	travisall = []string{
		"after_failure",
		"after_script",
		"after_success",
		"before_install",
		"before_script",
		"install",
		"script",
	}
	stages = []string{
		"before_install",
		"install",
		"before_script",
		"script",
		"after_script", // TODO(rjeczalik): #2
	}
)

type (
	travisProject  map[string]interface{}
	travisCommands map[string][]string
)

func travisParse(p []byte) (travisProject, error) {
	var pro = make(travisProject)
	if err := yaml.Unmarshal(p, pro); err != nil {
		return nil, err
	}
	for k, v := range pro {
		if v == nil {
			delete(pro, k)
			continue
		}
		i := sort.SearchStrings(travisall, k)
		if i >= len(travisall) || travisall[i] != k {
			delete(pro, k)
		}
	}
	return pro, nil
}

func travisCmd(pro travisProject) travisCommands {
	var m = make(travisCommands)
	for _, s := range travisall {
		v, ok := pro[s]
		if !ok {
			continue
		}
		if ok && v == nil {
			delete(pro, s)
			continue
		}
		switch cmd := v.(type) {
		case string:
			if cmd != "" {
				m[s] = []string{cmd}
			}
		case []interface{}:
			for _, cmd := range cmd {
				if cmd, ok := cmd.(string); ok && cmd != "" {
					m[s] = append(m[s], cmd)
				}
			}
		}
	}
	if len(m) != 0 {
		return m
	}
	return nil
}
