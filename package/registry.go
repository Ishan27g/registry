package _package

import (
	"sync"

	"github.com/Ishan27g/go-utils/mLogger"
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-hclog"
	"github.com/jedib0t/go-pretty/v6/table"
)

type MetaData interface{}
type peer RegisterRequest
type peers map[string]*peer // peer-address : peer
func (ps *peers) getPeers() []peer {
	var p []peer
	for _, p2 := range *ps {
		p = append(p, *p2)
	}
	return p
}

type registry struct {
	lock         sync.Mutex
	zones        *avltree.Tree // zoneId : peers
	logger       hclog.Logger
	serverEngine *gin.Engine
}

func (r *registry) clear() {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.zones.Clear()
}
func (r *registry) getPeers(zone int) []peer {
	r.lock.Lock()
	defer r.lock.Unlock()
	peersI, found := r.zones.Get(zone)
	if !found {
		return nil
	}
	peerMap := peersI.(peers)
	return peerMap.getPeers()
}

func (r *registry) addPeer(p peer) bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	added := false
	peersI, found := r.zones.Get(p.Zone)
	if !found {
		ps := peers{
			p.Address: &p,
		}
		r.zones.Put(p.Zone, ps)
		added = true
	} else {
		peerMap := peersI.(peers)
		if peerMap[p.Address] == nil { // new peer
			peerMap[p.Address] = &p
			added = true
		} else { // existing peer
			existingEntryForPeer := peerMap[p.Address]
			if existingEntryForPeer.RegisterAt.Before(p.RegisterAt) {
				peerMap[p.Address] = &p
				added = true
			}
		}
	}
	return added
}
func (r *registry) zoneIds() []int {
	r.lock.Lock()
	defer r.lock.Unlock()
	var zoneIds []int
	for _, i2 := range r.zones.Keys() {
		zoneIds = append(zoneIds, i2.(int))
	}
	return zoneIds
}
func (r *registry) allDetails(tbl bool) interface{} {
	r.lock.Lock()
	defer r.lock.Unlock()
	if !tbl {
		t := table.NewWriter()
		//t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.StyleLight)
		t.Style().Options.DrawBorder = false
		t.AppendHeader(table.Row{"Zone", "Peer Address", "Registered At"})

		for it := r.zones.Iterator(); it.Next(); {
			p := it.Value().(peers)
			var logs []table.Row
			for _, r := range p {
				logs = append(logs, table.Row{r.Zone, r.Address, r.RegisterAt.String()})
			}
			t.AppendRows(logs)
			t.AppendSeparator()
		}
		t.AppendSeparator()
		return t.Render()
	} else {
		allPeers := make(map[int]peers) // zone : []peers
		z := r.zones.Keys()
		for _, zone := range z {
			zoneId := zone.(int)
			p, _ := r.zones.Get(zoneId)
			allPeers[zoneId] = p.(peers)
		}
		return allPeers
	}

}

func Setup() *registry {
	reg := &registry{
		lock:   sync.Mutex{},
		zones:  avltree.NewWithIntComparator(),
		logger: mLogger.New("registry", "info"),
	}
	return reg
}
