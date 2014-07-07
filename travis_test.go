package ciexec

import (
	"reflect"
	"testing"
)

var travisFixture = [...]struct {
	cfg []byte
	pro travisProject
	cmd travisCommands
}{
	0: {
		[]byte(`
before_script: ./provision.sh
script:
 - go test ./...`),
		travisProject{
			"before_script": "./provision.sh",
			"script":        []interface{}{"go test ./..."},
		},
		travisCommands{
			"before_script": {
				"./provision.sh",
			},
			"script": {
				"go test ./...",
			},
		},
	},
	1: {
		[]byte(`
before_install: ./install.sh
install:
 - sudo apt-get update
 - sudo apt-get upgrade
before_script:
script:
 -
 - go test ./...
after_script:`),
		travisProject{
			"before_install": "./install.sh",
			"install": []interface{}{
				"sudo apt-get update",
				"sudo apt-get upgrade",
			},
			"script": []interface{}{nil, "go test ./..."},
		},
		travisCommands{
			"before_install": {
				"./install.sh",
			},
			"install": {
				"sudo apt-get update",
				"sudo apt-get upgrade",
			},
			"script": {
				"go test ./...",
			},
		},
	},
	2: {
		[]byte(`
before_install:
 - command 1
 - command 2
install:
 - command 5
 - command 6
 - command 7
before_script:
 - command 8
 - command 9
 - command 10
script: command 11
after_success:
 - command 12
 - command 13
after_failure: "./command 14"
after_script:
 - command 15
matrix:
 include:
  - 1
  - 2`),
		travisProject{
			"before_install": []interface{}{"command 1", "command 2"},
			"install":        []interface{}{"command 5", "command 6", "command 7"},
			"before_script":  []interface{}{"command 8", "command 9", "command 10"},
			"script":         "command 11",
			"after_success":  []interface{}{"command 12", "command 13"},
			"after_failure":  "./command 14",
			"after_script":   []interface{}{"command 15"},
		},
		travisCommands{
			"before_install": {
				"command 1",
				"command 2",
			},
			"install": {
				"command 5",
				"command 6",
				"command 7",
			},
			"before_script": {
				"command 8",
				"command 9",
				"command 10",
			},
			"script": {
				"command 11",
			},
			"after_success": {
				"command 12",
				"command 13",
			},
			"after_failure": {
				"./command 14",
			},
			"after_script": {
				"command 15",
			},
		},
	},
	3: {
		[]byte(`
language: go
go:
 - 1.2.2
 - 1.3
install:
 - go get code.google.com/p/go.tools/cmd/vet
 - go get -v ./...
 - go install -a -race std
script:
 - go tool vet -all .
 - go build ./...
 - go test -race -v ./...`),
		travisProject{
			"install": []interface{}{
				"go get code.google.com/p/go.tools/cmd/vet",
				"go get -v ./...",
				"go install -a -race std",
			},
			"script": []interface{}{
				"go tool vet -all .",
				"go build ./...",
				"go test -race -v ./...",
			},
		},
		travisCommands{
			"install": {
				"go get code.google.com/p/go.tools/cmd/vet",
				"go get -v ./...",
				"go install -a -race std",
			},
			"script": {
				"go tool vet -all .",
				"go build ./...",
				"go test -race -v ./...",
			},
		},
	},
}

func TestTravisParse(t *testing.T) {
	t.Parallel()
	for i, cas := range travisFixture {
		pro, err := travisParse(cas.cfg)
		if err != nil {
			t.Errorf("want err=nil; got %q (i=%d)", err, i)
			continue
		}
		if !reflect.DeepEqual(pro, cas.pro) {
			t.Errorf("want pro=%+v; got %+v (i=%d)", cas.pro, pro, i)
		}
	}
}

func TestTravisCmd(t *testing.T) {
	t.Parallel()
	for i, cas := range travisFixture {
		cmd := travisCmd(cas.pro)
		if !reflect.DeepEqual(cmd, cas.cmd) {
			t.Errorf("want cmd=%+v; got %+v (i=%d)", cas.cmd, cmd, i)
		}
	}
}
