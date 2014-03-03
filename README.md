coprocess
=========
[![GoDoc](https://godoc.org/github.com/kevinwallace/coprocess?status.png)](https://godoc.org/github.com/kevinwallace/coprocess)

`coprocess` allows a parent/child proccess pair to communicate using `net/rpc` semantics.  Communication is done over stdin/stdout of the child process, rather than over the network.

Usage
-----

On the parent:

    cmd := exec.Command("/path/to/child/binary", "--args", "--for", "--child")
    client, err := coprocess.NewClient(cmd)
    defer client.Close()
    // Use client to make RPCs to the child process...

On the child:

    s := rpc.NewServer()
    // Register RPC handlers on s...
    coprocess.Serve(s)

See the `examples` directory for a more complete example.
