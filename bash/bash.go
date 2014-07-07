// Package bash provides simple interface to a shell session.
package bash

import (
	"io"
	"os"
	"os/exec"

	"github.com/kr/pty"
)

var trap = []byte("__EXIT=0; trap '__EXIT=$?' ERR\n")
var exit = []byte("exit $__EXIT\n")

// Session represents a single Bash session. Any environment change are preserved
// across multiple commands.
type Session struct {
	w io.Writer
	d <-chan struct{}
	f *os.File
	c *exec.Cmd
}

// NewSession initiates new bash session. It sources user's default .bashrc file.
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

// Start runs a single line in the bash session. The execution is asynchronous,
// it does not wait for the command to finish.
func (s *Session) Start(cmd string) {
	s.f.Write([]byte(cmd))
	s.f.Write([]byte{'\n'})
}

// Close terminates the session. It waits for the session to cleanup. The function
// returns the last error encountered during the session, if any.
func (s *Session) Close() error {
	s.f.Write(exit)
	err := s.c.Wait()
	if w, ok := s.w.(io.Closer); ok {
		w.Close()
	}
	<-s.d
	return err
}
