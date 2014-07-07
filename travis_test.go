package ciexec

import (
	"fmt"
	"testing"
)

var output = []byte(`language: go
go:
 - 1.2.2
 - 1.3
install:
 - go get code.google.com/p/go.tools/cmd/vet
 - go get -v ./...
 - go  install -a -race std
script:
 - go tool vet -all .
 - go build ./...
 - go test -race -v ./...
`)

func TestTravisParse(t *testing.T) {
	t.Parallel()
	pro, err := travisParse(output)
	if err != nil {
		t.Errorf("want err=nil; got %q", err)
	}
	for k, v := range pro {
		fmt.Printf("%s: %v (%T)\n", k, v, v)
	}
}
