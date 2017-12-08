package main

import (
	"flag"
	"../../Gossiper/tools"
	"strings"
)

func main() {
	// Parsing inputs
	gossipPort := flag.String("gossipAddr", "localhost:5000", "gossipAddr:gossipPort")
	peers := flag.String("peers", "", "List of peers: addrPeer1:portPeer1_addrPeer2:portPeer2 ...")
	rtimer := flag.Uint("rtimer", 60, "Delay during two routing message (Developer)")
	flag.Parse()

	// Avoid :0 of being a peers if no intial peers are specified
	var peerAddrs []string
	if *peers != "" {
		peerAddrs = strings.Split(*peers, "_")
	} else {
		peerAddrs = make([]string, 0)
	}

	gossiper, _ := tools.NewGossiper("/tmp/gomat.sock", *gossipPort, "peerster", peerAddrs, *rtimer)
	daemon := Daemon{gossiper: gossiper}

	daemon.Run()
}
