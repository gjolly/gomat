package tools

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/matei13/gomat/Gossiper/tools/Peers"
)

type data struct {
	Capacity int      `json:"capacity"`
	Peers    []string `json:"peers"`
	Tasks    []Tasks  `json:"tasks"`
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
	peers := g.peers.Available()
	for _, p := peers {
		peersList = append(peersList, p.String())
	}
	
	jsonEncodeSend(w, data{
		Capacity: g.MaxCapacity,
		Peers:    peersList,
		Tasks:    g.vectorClock.GetMessages(),
	})
}

func (g *Gossiper) setCapacityHandler(w http.ResponseWriter, r *http.Request) {
	var capacity capacity
	json.NewDecoder(r.Body).Decode(&capacity)
	g.MaxCapacity = int(capacity)
}

func (g *Gossiper) addPeerHandler(w http.ResponseWriter, r *http.Request) {
	var peer peer
	json.NewDecoder(r.Body).Decode(&peer)
	addr, err := net.ResolveUDPAddr("udp4", string(peer))
	g.peers = append(g.peers, Peers.NewPeer(addr))
}

func (g *Gossiper) runServer(port string) {
	r := mux.NewRouter()
	r.HandleFunc("/getData", g.getDataHandler)
	r.HandleFunc("/setCapacity", g.setCapacityHandler)
	r.HandleFunc("/addPeer", g.addPeerHandler)
	http.ListenAndServe(":"+port, r)
}
