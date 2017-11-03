package main

import (
	"flag"
	"fmt"
	"strings"

	"./tools"
)

func main() {

	// Parsing inputs
	UIPort := flag.String("UIPort", "10000", "UIPort")
	gossipPort := flag.String("gossipPort", "localhost:5000", "gossipPort")
	nodeName := flag.String("name", "nodeA", "nodeName")
	peers := flag.String("peers", "", "peers")
	flag.Parse()

	// Avoid :0 of being a peers if no intial peers are specified
	var peerAddrs []string
	if *peers != "" {
		peerAddrs = strings.Split(*peers, "_")
	} else {
		peerAddrs = make([]string, 0)
	}

	// Creating a new gossiper
	g, err := tools.NewGossiper(*UIPort, *gossipPort, *nodeName, peerAddrs)
	if err != nil {
		fmt.Println(err)
		return
	}

	g.Run()

}
