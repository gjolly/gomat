package tools

import (
	"encoding/json"
	"net"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/matei13/gomat/Gossiper/tools/Peers"
	"github.com/matei13/gomat/Gossiper/tools/Pending"
)

type data struct {
	Capacity int            `json:"capacity"`
	Peers    []string       `json:"peers"`
	Info     []Pending.Info `json:"tasks"`
}

type capacity struct {
	Capacity int
}

type peer struct {
	Peer string
}

func jsonEncodeSend(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	json.NewEncoder(w).Encode(data)
}

func (g *Gossiper) getDataHandler(w http.ResponseWriter, r *http.Request) {
	var peersList []string
	peers := g.peers.Available(t1)
	for _, p := range peers {
		peersList = append(peersList, p.String())
	}

	jsonEncodeSend(w, data{
		Capacity: g.MaxCapacity,
		Peers:    peersList,
		Info:     g.Pending.GetInfos(),
	})
}

func (g *Gossiper) setCapacityHandler(w http.ResponseWriter, r *http.Request) {
	var capacity capacity
	json.NewDecoder(r.Body).Decode(&capacity)
	g.MaxCapacity = int(capacity.Capacity)
}

func (g *Gossiper) addPeerHandler(w http.ResponseWriter, r *http.Request) {
	var peer peer
	json.NewDecoder(r.Body).Decode(&peer)
	addr, _ := net.ResolveUDPAddr("udp4", string(peer.Peer))
	g.peers.Lock.Lock()
	defer g.peers.Lock.Unlock()
	g.peers.Map[addr.String()] = &Peers.Peer{Addr: *addr, Timer: 0}
}

func (g *Gossiper) RunServer(port string) {
	r := mux.NewRouter()
	r.HandleFunc("/getData", g.getDataHandler)
	r.HandleFunc("/setCapacity", g.setCapacityHandler)
	r.HandleFunc("/addPeer", g.addPeerHandler)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}
