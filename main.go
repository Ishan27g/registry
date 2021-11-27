package main

import (
	"os"

	"github.com/Ishan27gOrg/registry/package"
)

func main() {
	port := os.Getenv("BIND_ADDR")
	if port == "" {
		os.Exit(1)
	}
	_package.Server(port, _package.Setup())
}
