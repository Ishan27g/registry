package main

import (
	"os"

	_package "github.com/Ishan27gOrg/registry/package"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		os.Exit(1)
	}
	_package.Server(port, _package.Setup())
}
