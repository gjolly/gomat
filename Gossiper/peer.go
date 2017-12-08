package main

import (
	"flag"
	"fmt"
	"strings"
	"github.com/gomat/Gossiper/tools"
	"github.com/gomat/Gossiper/tools/Messages"
	"github.com/gomat/Gossiper/GUI"
)

type Peer struct {
	webServer *GUI.WebServer
	gossiper  *tools.Gossiper
}

func main() {

	// Parsing inputs
	UIPort := flag.String("UIPort", "10000", "UIPort")
	gossipPort := flag.String("gossipAddr", "localhost:5000", "gossipAddr:gossipPort")
	nodeName := flag.String("name", "nodeA", "Name of the node")
	peers := flag.String("peers", "", "List of peers: addrPeer1:portPeer1_addrPeer2:portPeer2 ...")
	rtimer := flag.Uint("rtimer", 60, "Delay during two routing message (Developer)")
	guiAddr := flag.String("guiAddr", "none", "Enable GUI. Address of the GUI have to be " +
		"precised: guiAddr:guiPort")
	flag.Parse()

	// Avoid :0 of being a peers if no intial peers are specified
	var peerAddrs []string
	if *peers != "" {
		peerAddrs = strings.Split(*peers, "_")
	} else {
		peerAddrs = make([]string, 0)
	}

	// Creating peer
	peer := Peer{}

	// Creating a new gossiper
	var err error
	peer.gossiper, err = tools.NewGossiper(*UIPort, *gossipPort, *nodeName, peerAddrs, *rtimer)
	if err != nil {
		fmt.Println(err)
		return
	}

	if *guiAddr != "none"{
		// Creating WebServer
		peer.webServer = GUI.NewWebServer(*guiAddr, peer.sendMsg, peer.sendPrivateMsg, &peer.gossiper.MessagesReceived,
			&peer.gossiper.PrivateMessages, &peer.gossiper.RoutingTable)
		fmt.Println("Peer: server addr=", peer.webServer.Addr)
		go 	peer.webServer.Run()
	}

	peer.gossiper.Run()
}

func (p *Peer) sendMsg(msg string) {
	p.gossiper.AcceptRumorMessage(Messages.RumorMessage{Text:msg}, *p.webServer.Addr, true)
}

func (p *Peer) sendPrivateMsg(msg, dest string){
	p.gossiper.AcceptRumorMessage(Messages.RumorMessage{Dest:dest, Text:msg}, *p.webServer.Addr, true)
}