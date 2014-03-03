package coprocess

import (
	"io"
	"net/rpc"
	"os"
	"os/exec"
)

// NewClient returns a new rpc.Client that communicates with the given command over stdin/stdout.
func NewClient(cmd *exec.Cmd) (*rpc.Client, error) {
	conn, err := newConn(cmd)
	if err != nil {
		return nil, err
	}
	return rpc.NewClient(conn), nil
}

// conn implements io.ReadWriteCloser,
// executing a command and reading/writing from its stdout/stdin.
type conn struct {
	cmd    *exec.Cmd
	stdout io.ReadCloser
	stdin  io.WriteCloser
	// When the process exits, err is set to the result of cmd.Wait, and done is closed.
	done chan bool
	err  error
}

// newConn executes the given command, returning an io.ReadWriteCloser attached to its stdout/stdin.
// The Cmd's Stdout/Stdin/Stderr are modified, but everything else (env, cwd, etc) is left intact.
// Takes ownership of cmd -- callers may not modify it after passing it to NewConn.
func newConn(cmd *exec.Cmd) (io.ReadWriteCloser, error) {
	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	conn := &conn{
		cmd:    cmd,
		stdout: stdout,
		stdin:  stdin,
		done:   make(chan bool),
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	go func() {
		conn.err = cmd.Wait()
		close(conn.done)
	}()
	return conn, nil
}

func (c *conn) wrapError(err error) error {
	if err == io.EOF {
		<-c.done
		if c.err != nil {
			return c.err
		}
	}
	return err
}

func (c *conn) Read(p []byte) (int, error) {
	n, err := c.stdout.Read(p)
	return n, c.wrapError(err)
}

func (c *conn) Write(p []byte) (int, error) {
	n, err := c.stdin.Write(p)
	return n, c.wrapError(err)
}

func (c *conn) Close() error {
	if err := c.stdin.Close(); err != nil {
		return err
	}
	<-c.done
	return c.err
}
