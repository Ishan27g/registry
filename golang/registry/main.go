package main

import (
	"os"

	p "github.com/Ishan27g/registry/golang/registry/package"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		os.Exit(1)
	}
	p.Run(port, p.Setup())

	//c := _package.NewClient("http://localhost:9999")
	//fmt.Println(c.Register("http://localhost:1000", map[string]string{
	//	"ok": "okkok",
	//}))
	////go fmt.Println(c.GetZonePeers(1))
	////go fmt.Println(c.GetZonePeers(2))
	////fmt.Println(c.GetZonePeers(3))
	//
	//fmt.Println(c.GetZoneIds())
	//
}
