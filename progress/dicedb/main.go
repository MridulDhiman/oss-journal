//go:build linux
// +build linux

// main.go
package main

import (
	"flag"

	"github.com/MridulDhiman/dice/config"
	"github.com/MridulDhiman/dice/server"
)

func init() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host for the dice server")
	flag.IntVar(&config.Port, "port", 7379, "port no. for the dice server")
	// parses the flags from actual CLI flags
	flag.Parse() // it need to be called, after the flags are defined and before they are used
}

func main() {
	server.RunAsyncTCPServer()
}
