package main

import (
	"log"
	"net/rpc"

	"github.com/kevinwallace/coprocess"
)

type Args struct {
	A, B int
}

type Arith struct{}

func (a *Arith) Add(args *Args, reply *int) error {
	log.Printf("[server] got request: %#v", args)
	*reply = args.A + args.B
	log.Printf("[server] returning: %d", *reply)
	return nil
}

func main() {
	defer log.Print("[server] shutting down...")
	s := rpc.NewServer()
	s.Register(&Arith{})
	coprocess.Serve(s)
}
