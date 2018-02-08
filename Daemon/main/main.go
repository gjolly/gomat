package main

import (
	"flag"
	"github.com/matei13/gomat/Gossiper/tools"
	"strings"
	"github.com/matei13/gomat/Daemon"
)

func main() {
	// Parsing inputs
	gossipPort := flag.String("gossipAddr", "localhost:5000", "gossipAddr:gossipPort")
	peers := flag.String("peers", "", "List of peers: addrPeer1:portPeer1_addrPeer2:portPeer2 ...")
	capacity := flag.Int("capacity", 5000, "Computing power of the gossiper")
	flag.Parse()

	// Avoid :0 of being a peers if no intial peers are specified
	var peerAddrs []string
	if *peers != "" {
		peerAddrs = strings.Split(*peers, "_")
	} else {
		peerAddrs = make([]string, 0)
	}

	gossiper, err := tools.NewGossiper(*gossipPort, "peerster", peerAddrs, *capacity)
	if err != nil {
		panic(err)
	}
	daemon := Daemon.Daemon{Gossiper: gossiper}
	daemon.Run()
}
