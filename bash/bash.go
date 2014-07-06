package bash

import (
	"io"
	"os"
	"os/exec"

	"github.com/kr/pty"
)

var trap = []byte("__EXIT=0; trap '__EXIT=$?' ERR\n")
var exit = []byte("exit $__EXIT\n")

// Session TODO
type Session struct {
	w io.Writer
	d <-chan struct{}
	f *os.File
	c *exec.Cmd
}

// NewSession TODO
func NewSession(w io.Writer) (*Session, error) {
	c := exec.Command("bash")
	c.Args = []string{"-i"}
	f, err := pty.Start(c)
	if err != nil {
		return nil, err
	}
	d := make(chan struct{})
	go func() {
		io.Copy(w, f)
		close(d)
	}()
	s := &Session{w: w, d: d, f: f, c: c}
	if _, err = s.f.Write(trap); err != nil {
		s.Close()
		return nil, err
	}
	return s, nil
}

// Start TODO
func (s *Session) Start(cmd string) {
	s.f.Write([]byte(cmd))
	s.f.Write([]byte{'\n'})
}

// Close TODO
func (s *Session) Close() error {
	s.f.Write(exit)
	err := s.c.Wait()
	if w, ok := s.w.(io.Closer); ok {
		w.Close()
	}
	<-s.d
	return err
}
