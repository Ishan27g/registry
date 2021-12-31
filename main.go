package main

import (
	"os"

	p "github.com/Ishan27gOrg/registry/package"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		os.Exit(1)
	}
	p.Run(port, p.Setup())
}
