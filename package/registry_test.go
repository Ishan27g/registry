package _package

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testAddr = "9999"

const testUrl = "http://localhost"

// const deployedUrl = "https://bootstrap-registry.herokuapp.com"
const tClientsPerZone = 5
const tZones = 6
const MockHostName = "http://localhost:"

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
		addr := MockHostName + strconv.Itoa(port)
		nw[addr] = newPeer(addr, zone, nil)
	}
	return nw
}
func mockClientGetZones() []int {
	// return RegistryClient(deployedUrl).GetZoneIds()
	return RegistryClient(testUrl + ":" + testAddr).GetZoneIds()
}
func mockClientGetZone(zone int) PeerResponse {
	// return RegistryClient(deployedUrl).GetZonePeers(zone)
	return RegistryClient(testUrl + ":" + testAddr).GetZonePeers(zone)
}
func mockClientRegister(zone int, address string, meta MetaData) PeerResponse {
	//return RegistryClient(deployedUrl).Register(zone, address, meta)
	return RegistryClient(testUrl+":"+testAddr).Register(zone, address, meta)
}
func mockClientGetDetails() []string {
	return RegistryClient(testUrl + ":" + testAddr).GetDetails()
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

	assert.Equal(t, len(mockClientGetZones()), tZones)

	for i, zone := range testNetwork {
		for _, p := range zone {
			check[i] = mockClientGetZone(p.Zone)
		}
	}
	wg2.Wait()

	assert.NotEmpty(t, mockClientGetDetails())

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

	fmt.Printf("\n%v\n", mockClientGetDetails()) // = reg.allDetails(true)
	fmt.Println(reg.allDetails(false))
}
