package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	RegUrl       = "/register"
	ZoneIdsUrl   = "/zones"
	ZonePeersUrl = "/zone"
	//DetailsUrl     = "/details"
	DetailsUrlJson = "/details/json"
	ShutdownUrl    = "/shutdown"
	ResetUrl       = "/reset"
)

type Client interface {
	// Register self at this address/zone with registry
	Register(address string, meta MetaData) PeerResponse
	// GetZoneIds returns the zoneIds
	GetZoneIds() []int
	// GetZonePeers returns the addresses of zone peers
	GetZonePeers(zone int) PeerResponse

	ping(address string) bool
}

type MetaData interface{}
type jsonRequest struct {
	RegisterAt time.Time `json:"registered_at"`
	Address    string    `json:"address"` // full Address `
	Zone       int       `json:"zone"`    // todo
	MetaData   MetaData  `json:"meta_data"`
}
type PeerResponse []jsonRequest

func (c *client) Register(address string, meta MetaData) PeerResponse {
	if meta == nil {
		return nil
	}
	type reqJson struct {
		Address  string   `json:"address"`
		MetaData MetaData `json:"metaData"`
	}
	var r reqJson
	r.Address = address
	r.MetaData = meta
	json, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	req, err := http.NewRequest("POST", c.serverAddress+RegUrl, bytes.NewBuffer(json))
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	return parseRegRsp(sendReq(req))
}

func (c *client) GetZoneIds() []int {
	url := c.serverAddress + ZoneIdsUrl
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}
	var rsp []int
	if b := sendReq(req); b != nil {
		err := json.Unmarshal(b, &rsp)
		if err != nil {
			return nil
		}
	}
	return rsp
}

func (c *client) GetZonePeers(zone int) PeerResponse {
	url := c.serverAddress + ZonePeersUrl + "/" + strconv.Itoa(zone)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}
	return parseRegRsp(sendReq(req))
}
func parseRegRsp(j []byte) []jsonRequest {
	if j != nil {
		var rr []jsonRequest
		err := json.Unmarshal(j, &rr)
		if err != nil {
			fmt.Println("e" + err.Error())
		}
		return rr
	}
	return nil
}
func sendReq(req *http.Request) []byte {
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR reading response " + err.Error())
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR reading body. " + err.Error())
		return nil
	}
	defer resp.Body.Close()
	return body
}
func (c *client) ping(address string) bool {
	req, err := http.NewRequest("GET", address+"/engine/ping", nil)
	if err != nil {
		return false
	}
	client := &http.Client{Timeout: time.Second * 3}

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	return resp.StatusCode == http.StatusOK
}

type client struct {
	serverAddress string
}

func NewClient(serverAddress string) Client {
	return &client{serverAddress: serverAddress}
}
