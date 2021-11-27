package _package

import (
	"context"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testAddr = ":9999"
const testUrl = "http://localhost"
const tClientsPerZone = 3
const tZones = 5

var reg *registry
var singleRegistry = sync.Once{}
var ctx context.Context
var cancel context.CancelFunc
var stopRegistry chan bool

func mockRegistry() {
	singleRegistry.Do(func() {
		stopRegistry = make(chan bool)
		ctx, cancel = context.WithCancel(context.Background())
		reg = Setup()
		go Server(testAddr, reg)
		go func() {
			<-stopRegistry
			cancel()
		}()

	})
	<-time.After(300 * time.Millisecond)
}

type testZone map[string]peer

func mockZone(zone int) testZone {
	nw := make(map[string]peer)
	for i := 0; i < tClientsPerZone; i++ {
		port := (zone * 10) + 9000 + i
		addr := "http://service:" + strconv.Itoa(port)
		nw[addr] = newPeer(addr, zone, nil)
	}
	return nw
}
func mockClientGetZone(zone int) PeerResponse {
	return RegistryClient(testUrl + testAddr).GetZonePeers(zone)
}
func mockClientRegister(zone int, address string, meta MetaData) PeerResponse {
	return RegistryClient(testUrl+testAddr).Register(zone, address, meta)
}

func TestRegistryClient(t *testing.T) {
	mockRegistry()
	rand.Seed(time.Now().Unix())

	testNetwork := make(map[int]testZone)
	for i := 0; i < tZones; i++ {
		testNetwork[i+1] = mockZone(i + 1)
	}
	wg := sync.WaitGroup{}
	rsp1 := make(map[string]PeerResponse)
	mapLock := sync.Mutex{}
	for _, zone := range testNetwork {
		for _, p := range zone {
			wg.Add(1)
			go func(p peer, wg *sync.WaitGroup) {
				defer wg.Done()
				mapLock.Lock()
				rsp1[p.Address] = mockClientRegister(p.Zone, p.Address, p.MetaData)
				mapLock.Unlock()
				x := time.Duration(rand.Intn(1000))
				time.Sleep(time.Millisecond * x) // re-register again after random failure/timeout
				mapLock.Lock()
				rsp1[p.Address] = mockClientRegister(p.Zone, p.Address, p.MetaData)
				mapLock.Unlock()
			}(p, &wg)
		}
	}
	wg.Wait()

	wg2 := sync.WaitGroup{}
	check := make(map[int]PeerResponse)
	for i, zone := range testNetwork {
		for _, p := range zone {
			check[i] = mockClientGetZone(p.Zone)
		}
	}
	wg2.Wait()

	for zoneNum, responses := range check {
		// DEBUG PRINTS
		// str := fmt.Sprintf("\n============= Zone %d ============", zoneNum)
		// for _, r := range responses {
		// str +=  fmt.Sprintf("\n [Zone %d] - %s (%s) %v", r.Zone, r.Address, r.RegisterAt.String(), r.MetaData)
		// }
		//	fmt.Println(str)
		for _, zonePeer := range responses {
			assert.Equal(t, zoneNum, zonePeer.Zone)
		}
		assert.Equal(t, tClientsPerZone, len(responses))
	}

	reg.allDetails()

}
