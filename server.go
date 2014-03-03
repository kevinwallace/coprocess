package coprocess

import (
	"io"
	"net/rpc"
	"os"
)

type stdio struct {
	stdin  io.ReadCloser
	stdout io.WriteCloser
}

func (s stdio) Read(p []byte) (int, error) {
	return s.stdin.Read(p)
}

func (s stdio) Write(p []byte) (int, error) {
	return s.stdout.Write(p)
}

func (s stdio) Close() error {
	if err := s.stdout.Close(); err != nil {
		return err
	}
	if err := s.stdin.Close(); err != nil {
		return err
	}
	return nil
}

// Serve handles RPC commands from stdin, writing responses to stdout, until stdin is closed.
func Serve(server *rpc.Server) {
	server.ServeConn(stdio{os.Stdin, os.Stdout})
}
