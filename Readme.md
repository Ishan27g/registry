```go
package main

import reg "github.com/Ishan27gOrg/registry/package"

type RegistryClientI interface {
	// Register self at this address/zone with registry
	Register(zone int, address string, meta reg.MetaData) reg.PeerResponse
	// GetZoneIds returns the zoneIds
	GetZoneIds() []int
	// GetZonePeers returns the addresses of zone peers
	GetZonePeers(zone int) reg.PeerResponse
	// GetDetails returns all registered peers details
	GetDetails() []string
}
```
