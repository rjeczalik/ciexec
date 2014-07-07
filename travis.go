package ciexec

import (
	"errors"
	"io"
	"path/filepath"

	"gopkg.in/yaml.v1"
)

func init() {
	all = append(all, travis{})
}

type travis struct{}

func (travis) Is(file string, _ []byte) bool {
	return filepath.Base(file) == ".travis.yml"
}

func (travis) Exec(stage string, p []byte, w io.Writer) error {
	return nil
}

var order = []string{
	"before_install",
	"install",
	"before_script",
	"script",
}

// TravisProject TODO
type TravisProject map[string]interface{}

func travisParse(p []byte) (TravisProject, error) {
	var pro = make(TravisProject)
	if err := yaml.Unmarshal(p, pro); err != nil {
		return nil, err
	}
	var n int
	for _, s := range order {
		if cmd, ok := pro[s]; ok {
			if cmd == nil {
				delete(pro, s)
			} else {
				n += 1
			}
		}
	}
	if n == 0 {
		return nil, errors.New("no stages found or they have no commands")
	}
	return pro, nil
}
