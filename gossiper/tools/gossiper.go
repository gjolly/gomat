package tools

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/dedis/protobuf"
)

// Gossiper -- Discripe a node of a Gossip network
type Gossiper struct {
	UIAddr           *(net.UDPAddr)
	gossipAddr       *(net.UDPAddr)
	UIConn           *(net.UDPConn)
	gossipConn       *(net.UDPConn)
	name             string
	peers            map[string]Peer
	vectorClock      []PeerStatus
	idMessage        int
	messagesReceived map[string](map[int]RumorMessage)
	exchangeEnded    chan bool
}

// NewGossiper -- Return a new gossiper structure
func NewGossiper(UIPort, gossipPort, identifier string, peerAddrs []string) (*Gossiper, error) {
	// For UIPort
	UIUdpAddr, err := net.ResolveUDPAddr("udp4", ":"+UIPort)
	if err != nil {
		return nil, err
	}
	UIConn, err := net.ListenUDP("udp4", UIUdpAddr)
	if err != nil {
		return nil, err
	}
	// For gossipPort
	gossipUdpAddr, err := net.ResolveUDPAddr("udp4", gossipPort)
	if err != nil {
		return nil, err
	}
	gossipConn, err := net.ListenUDP("udp4", gossipUdpAddr)
	if err != nil {
		return nil, err
	}

	g := &Gossiper{
		UIAddr:           UIUdpAddr,
		gossipAddr:       gossipUdpAddr,
		UIConn:           UIConn,
		gossipConn:       gossipConn,
		name:             identifier,
		peers:            make(map[string]Peer, 0),
		vectorClock:      make([]PeerStatus, 0),
		idMessage:        1,
		messagesReceived: make(map[string](map[int]RumorMessage), 0),
		exchangeEnded:    make(chan bool),
	}

	for _, peerAddr := range peerAddrs {
		udpAddr, err := net.ResolveUDPAddr("udp4", peerAddr)
		g.AddPeer(*udpAddr)
		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
	}
	return g, nil
}

func (g Gossiper) excludeAddr(excludedAddrs string) (addrs []string) {
	addrs = make([]string, 0)
	for addrPeer := range g.peers {
		if addrPeer != excludedAddrs {
			addrs = append(addrs, addrPeer)
		}
	}
	return
}

func (g Gossiper) getRandomPeer(excludedAddrs string) *Peer {
	availableAddrs := g.excludeAddr(excludedAddrs)
	//fmt.Println(availableAddrs)

	if len(availableAddrs) == 0 {
		return nil
	}

	i := rand.Intn(len(availableAddrs))

	peer := g.peers[availableAddrs[i]]
	return &peer
}

// AddPeer -- Add a new peer to the list of peers. If Peer is already known: do nothing
func (g *Gossiper) AddPeer(address net.UDPAddr) {
	IPAddress := address.String()
	_, okAddr := g.peers[IPAddress]
	if !okAddr {
		g.peers[IPAddress] = newPeer(address)
	}
}

func (g Gossiper) sendMessage(mess GossipMessage, addr net.UDPAddr) error {
	messEncode, err := protobuf.Encode(&mess)
	if err != nil {
		return err
	}
	g.gossipConn.WriteToUDP(messEncode, &addr)
	return nil
}

//send a message to all known peers excepted Peer
func (g Gossiper) sendToAllPeers(mess GossipMessage, excludeAddr *net.UDPAddr) {
	for _, p := range g.peers {
		if p.addr.String() != excludeAddr.String() {
			g.sendMessage(mess, p.addr)
		}
	}
}

func (g *Gossiper) listenConn(conn *net.UDPConn, isFromClient bool) {
	var bufferMess []byte
	var nbBytes int
	var err error
	var addr *net.UDPAddr
	for {
		bufferMess = make([]byte, 2048)
		nbBytes, addr, err = conn.ReadFromUDP(bufferMess)
		if err == nil {
			go g.accept(bufferMess, addr, nbBytes, isFromClient)
		}
	}
}

// Run -- Launch the server
func (g *Gossiper) Run() {
	go g.antiEntropy()
	go g.listenConn(g.UIConn, true)
	g.listenConn(g.gossipConn, false)
}

func (g *Gossiper) accept(buffer []byte, addr *net.UDPAddr, nbByte int, isFromClient bool) {
	mess := &GossipMessage{}
	protobuf.Decode(buffer, mess)

	if mess.Rumor != nil {
		g.acceptRumorMessage(*mess.Rumor, *addr, isFromClient)
	} else if mess.Status != nil {
		g.acceptStatusMessage(*mess.Status, addr)
	}
}

//Callback function, call when a message is received
func (g *Gossiper) acceptRumorMessage(mess RumorMessage, addr net.UDPAddr, isFromClient bool) {

	if !isFromClient && g.alreadySeen(mess.PeerMessage.ID, mess.Origin) {
		return
	}

	g.printDebugRumor(mess, mess.Origin, addr.String(), isFromClient)

	if isFromClient {
		mess.Origin = g.name
		mess.PeerMessage.ID = g.idMessage
		g.idMessage++
	}

	g.updateVectorClock(mess.Origin, mess.PeerMessage.ID)
	g.storeRumorMessage(mess, mess.PeerMessage.ID, mess.Origin)

	if !isFromClient {
		g.sendStatusMessage(addr)
		g.AddPeer(addr)
	}

	g.propagateRumorMessage(mess, addr.String())
}

func (g Gossiper) propagateRumorMessage(mess RumorMessage, excludedAddrs string) {
	coin := 1
	peer := g.getRandomPeer(excludedAddrs)

	for coin == 1 && peer != nil {

		fmt.Println("MONGERING with", peer.addr.String())
		g.sendMessage(GossipMessage{Rumor: &mess}, peer.addr)

		peer = g.getRandomPeer("")
		coin = rand.Int() % 2
		//fmt.Println(coin, peer)
		if coin == 1 && peer != nil {
			fmt.Println("FLIPPED COIN sending rumor to", peer.addr.String())
		}
	}
}

func (g Gossiper) acceptStatusMessage(mess StatusMessage, addr *net.UDPAddr) {
	messToSend, isMessageToAsk := g.compareVectorClocks(mess.Want)
	g.printDebugStatus(mess, *addr)

	if messToSend != nil {
		fmt.Println("MONGERING with", addr.String())
		g.sendMessage(*messToSend, *addr)
	}
	if isMessageToAsk {
		g.sendStatusMessage(*addr)
	}
	if !isMessageToAsk && messToSend == nil {
		fmt.Println("IN SYNC WITH", addr.String())
		g.exchangeEnded <- true
	}

	g.AddPeer(*addr)
}

func (g Gossiper) printPeerList() {
	first := true
	for _, peer := range g.peers {
		if first {
			first = false
			fmt.Print(peer)
		} else {
			fmt.Print(",", peer)
		}
	}
	fmt.Println()
}

func (g Gossiper) printDebugStatus(mess StatusMessage, addr net.UDPAddr) {
	fmt.Print("STATUS from ", addr.String())
	for _, peerStatus := range mess.Want {
		fmt.Print(" origin ", peerStatus.Identifier, " nextID ", peerStatus.NextID)
	}
	fmt.Println()
	g.printPeerList()
}

func (g Gossiper) printDebugRumor(mess RumorMessage, emitterName, lastHopIP string, isFromClient bool) {
	if isFromClient {
		fmt.Println("CLIENT", mess, g.name)
	} else {
		fmt.Println("RUMOR", "origin", emitterName, "from", lastHopIP, "ID", mess.PeerMessage.ID, "contents", mess.PeerMessage.Text)
	}
	g.printPeerList()
}

func (g Gossiper) alreadySeen(id int, nodeName string) bool {
	_, ok := g.messagesReceived[nodeName][id]
	return ok
}

func (g *Gossiper) updateVectorClock(name string, id int) {
	find := false
	for i := range g.vectorClock {
		if g.vectorClock[i].Identifier == name {
			find = true
			if g.vectorClock[i].NextID == id {
				g.vectorClock[i].NextID++
				g.checkOnAlreadySeen(g.vectorClock[i].NextID, name)
			}
		}
	}
	if !find && id == 1 {
		g.vectorClock = append(g.vectorClock, PeerStatus{Identifier: name, NextID: 2})
		g.checkOnAlreadySeen(2, name)
	}
}

func (g Gossiper) checkOnAlreadySeen(nextID int, nodeName string) {
	if g.alreadySeen(nextID, nodeName) {
		g.updateVectorClock(nodeName, nextID)
	}
}

func (g Gossiper) storeRumorMessage(mess RumorMessage, id int, nodeName string) {
	if g.messagesReceived[nodeName] == nil {
		g.messagesReceived[nodeName] = make(map[int]RumorMessage)
	}
	g.messagesReceived[nodeName][id] = mess
}

func (g *Gossiper) sendStatusMessage(addr net.UDPAddr) {
	g.sendMessage(GossipMessage{Status: &StatusMessage{Want: g.vectorClock}}, addr)
}

func (g *Gossiper) antiEntropy() {
	tick := time.NewTicker(time.Second)
	for {
		<-tick.C
		peer := g.getRandomPeer("")
		if peer != nil {
			g.sendStatusMessage(peer.addr)
		}
	}
}

func (g Gossiper) compareVectorClocks(vectorClock []PeerStatus) (msgToSend *GossipMessage, isMessageToAsk bool) {
	isMessageToAsk = false
	msgToSend = nil
	var find bool
	for _, ps := range g.vectorClock {
		find = false
		for _, wanted := range vectorClock {
			if ps.Identifier == wanted.Identifier {
				find = true
				if ps.NextID < wanted.NextID && wanted.NextID > 0 {
					isMessageToAsk = true
					if msgToSend != nil {
						return
					}
				} else if ps.NextID > wanted.NextID && wanted.NextID > 0 {
					mess := g.messagesReceived[ps.Identifier][ps.NextID]
					msgToSend = &GossipMessage{Rumor: &mess}
					if isMessageToAsk {
						return
					}
				}
			}
		}
		if !find {
			rm := g.messagesReceived[ps.Identifier][0]
			msgToSend = &GossipMessage{Rumor: &rm}
			if isMessageToAsk {
				return
			}
		}
	}
	for _, wanted := range vectorClock {
		find = false
		for _, ps := range g.vectorClock {
			if ps.Identifier == wanted.Identifier {
				find = true
			}
		}
	}
	if !find {
		isMessageToAsk = true
		return
	}
	return
}
