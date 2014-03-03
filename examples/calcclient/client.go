package main

import (
	"log"
	"os/exec"

	"github.com/kevinwallace/coprocess"
)

type Args struct {
	A, B int
}

func main() {
	client, err := coprocess.NewClient(exec.Command("calcserver"))
	defer client.Close()
	if err != nil {
		panic(err)
	}
	add := func(a, b int) (result int) {
		if err := client.Call("Arith.Add", Args{a, b}, &result); err != nil {
			panic(err)
		}
		return
	}
	log.Printf("[client] add(1, 2) = %d", add(1, 2))
	log.Printf("[client] add(3, 4) = %d", add(3, 4))
}
