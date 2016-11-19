package main

import (
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage %s
	server - run as server
	client - run as client
`,
		os.Args[0])
}
func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	switch os.Args[1] {
	case "server":
		server(os.Args[1:])
	case "client":
		client(os.Args[1:])
	default:
		usage()
	}
}
