package tools

import (
	"fmt"
	"math/rand"
	"net"
	"time"
	"github.com/dedis/protobuf"
	"sync"
	"github.com/matei13/gomat/Gossiper/tools/Messages"
	"github.com/matei13/gomat/Daemon/gomatcore"
	"github.com/matei13/gomat/Gossiper/tools/Peers"
	"github.com/matei13/gomat/Gossiper/tools/Tasks"
)

// Gossiper -- Describe a node of a Gossip network
type Gossiper struct {
	UIAddr           *net.UnixAddr
	gossipAddr       *net.UDPAddr
	UIListener       *net.UnixListener
	gossipConn       *net.UDPConn
	name             string
	peers            Peers.PeerMap
	idMessage        uint32
	MessagesReceived map[string]map[uint32]Messages.RumorMessage
	exchangeEnded    chan bool
	RoutingTable     RoutingTable
	mutex            *sync.Mutex
	rtimer           uint
	PrivateMessages  []Messages.RumorMessage
	MaxCapacity      int
	CurrentCapacity  int
	Tasks            Tasks.TaskMap
	Pending          map[string]map[int]chan bool
	OwnTasks         map[int]map[int]string
	TaskId           int
}

const t1 = 2

const timer = 30

const buffSize = 65507

// NewGossiper -- Returns a new gossiper structure
func NewGossiper(sockFile, gossipPort, identifier string, peerAddrs []string, rtimer uint, capa int) (*Gossiper, error) {
	// For UIPort
	UIAddr, err := net.ResolveUnixAddr("unix", sockFile)
	if err != nil {
		return nil, err
	}
	UIListener, err := net.ListenUnix("unix", UIAddr)
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
		UIAddr:           UIAddr,
		gossipAddr:       gossipUdpAddr,
		UIListener:       UIListener,
		gossipConn:       gossipConn,
		name:             identifier,
		peers:            Peers.PeerMap{Map: make(map[string] *Peers.Peer), Lock: &sync.RWMutex{}},
		idMessage:        1,
		MessagesReceived: make(map[string]map[uint32]Messages.RumorMessage, 0),
		exchangeEnded:    make(chan bool),
		RoutingTable:     *newRoutingTable(),
		mutex:            &sync.Mutex{},
		rtimer:           rtimer,
		MaxCapacity:      capa,
		CurrentCapacity:  capa,
		Tasks:            Tasks.TaskMap{Tasks: make(map[string][]Tasks.Task), Lock: &sync.RWMutex{}},
		Pending:          make(map[string]map[int]chan bool),
		OwnTasks:         make(map[int]map[int]string),
		TaskId:           0,
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
	g.peers.Lock.Lock()
	defer g.peers.Lock.Unlock()
	for addrPeer := range g.peers.Map {
		if addrPeer != excludedAddrs {
			addrs = append(addrs, addrPeer)
		}
	}
	return
}

func (g Gossiper) getRandomPeer(excludedAddrs string) *Peers.Peer {
	availableAddrs := g.excludeAddr(excludedAddrs)
	//fmt.Println(availableAddrs)

	if len(availableAddrs) == 0 {
		return nil
	}

	i := rand.Intn(len(availableAddrs))

	peer, _ := g.peers.Get(availableAddrs[i])
	return peer
}

// AddPeer -- Adds a new peer to the list of peers. If Peer is already known: do nothing
func (g *Gossiper) AddPeer(address net.UDPAddr) {
	IPAddress := address.String()
	g.peers.Lock.Lock()
	defer g.peers.Lock.Unlock()
	_, okAddr := g.peers.Map[IPAddress]
	if !okAddr {
		g.peers.Map[IPAddress] = &Peers.Peer{Addr: address, Timer: 0}
	}
	go g.listenGossiper(address)
}

func (g Gossiper) sendRumorMessage(message Messages.RumorMessage, addr net.UDPAddr) error {
	gossipMessage := Messages.GossipMessage{Rumor: &message}
	messEncode, err := protobuf.Encode(&gossipMessage)
	if err != nil {
		fmt.Println("error protobuf")
		return err
	}
	g.gossipConn.WriteToUDP(messEncode, &addr)
	return nil
}

func (g Gossiper) sendStatusMessage(addr net.UDPAddr) error {
	messEncode, err := protobuf.Encode(&Messages.GossipMessage{Status: &Messages.StatusMessage{}})
	if err != nil {
		fmt.Println("error protobuf")
		return err
	}
	g.gossipConn.WriteToUDP(messEncode, &addr)
	return nil
}

func (g *Gossiper) listenConn(conn *net.UDPConn) {
	var bufferMess []byte
	var nbBytes int
	var err error
	var addr *net.UDPAddr
	for {
		bufferMess = make([]byte, buffSize)
		nbBytes, addr, err = conn.ReadFromUDP(bufferMess)
		if err == nil {
			go g.accept(bufferMess, addr, nbBytes, false)
		}
	}
}

func (g *Gossiper) listenUnix(listener *net.UnixListener) {
	for {
		unixConn, _ := listener.AcceptUnix()
		g.acceptUI(unixConn)
	}
}

func (g *Gossiper) acceptUI(conn *net.UnixConn) {
	var bufferMess []byte
	var nbBytes int
	var err error
	for {
		bufferMess = make([]byte, buffSize)
		nbBytes, _, err = conn.ReadFromUnix(bufferMess)
		if err == nil {
			go g.accept(bufferMess, nil, nbBytes, true)
		}
	}
}

// Run -- Launch the server
func (g *Gossiper) Run() {
	go g.listenUnix(g.UIListener)
	go g.listenConn(g.gossipConn)
	go g.antiEntropy()
	g.sendRouteRumor()
	g.routeRumorDaemon()
}

func (g *Gossiper) accept(buffer []byte, addr *net.UDPAddr, nbByte int, isFromClient bool) {
	mess := &Messages.GossipMessage{}
	protobuf.Decode(buffer, mess)
	if mess.Rumor != nil {
		g.AcceptRumorMessage(*mess.Rumor, *addr, isFromClient)
	} else if mess.Status != nil {
		g.acceptStatusMessage(*mess.Status, addr)
	}
}

//Callback function, call when a message is received
func (g *Gossiper) AcceptRumorMessage(mess Messages.RumorMessage, addr net.UDPAddr, isFromClient bool) {

	if !isFromClient && g.alreadySeen(mess.ID, mess.Origin) {
		return
	}

	if isFromClient {
		mess.Origin = g.name
		if !mess.IsPrivate() {
			g.mutex.Lock()
			mess.ID = g.idMessage
			g.idMessage++
			g.mutex.Unlock()
		} else {
			mess.HopLimit = 10
		}
	} else {
		g.mutex.Lock()
		fmt.Println("DSDV", mess.Origin+":"+addr.String())
		g.RoutingTable.add(mess.Origin, addr.String())
		g.mutex.Unlock()
	}

	if mess.Text != "" {
		g.printDebugRumor(mess, addr.String(), isFromClient)
	}

	if !mess.IsPrivate() {
		g.storeRumorMessage(mess, mess.ID, mess.Origin)

		if !isFromClient {
			g.sendStatusMessage(addr)
			g.AddPeer(addr)
		}
		g.propagateRumorMessage(mess, addr.String())
	} else {
		if mess.HopLimit > 1 && mess.Dest != g.name {
			g.forward(mess)
		} else if mess.Dest == g.name {
			g.receivePrivateMessage(mess)
		}
	}
}

func (g *Gossiper) receivePrivateMessage(message Messages.RumorMessage) {
	fmt.Println("PRIVATE:", message.Origin+":"+fmt.Sprint(message.HopLimit)+":"+message.Text)
	g.PrivateMessages = append(g.PrivateMessages, message)
}

func (g Gossiper) forward(message Messages.RumorMessage) {
	message.HopLimit -= 1
	addr := g.RoutingTable.FindNextHop(message.Dest)
	if addr != "" {
		UDPAddr, err := net.ResolveUDPAddr("udp4", addr)
		if err == nil {
			fmt.Println("FORWARD private msg", message.Dest, addr)
			g.sendRumorMessage(message, *UDPAddr)
		}
	}
}

func (g Gossiper) propagateRumorMessage(mess Messages.RumorMessage, excludedAddrs string) {
	coin := 1
	peer := g.getRandomPeer(excludedAddrs)

	for coin == 1 && peer != nil {

		fmt.Println("MONGERING with", peer.Addr.String())
		g.sendRumorMessage(mess, peer.Addr)

		peer = g.getRandomPeer("")
		coin = rand.Int() % 2
		//fmt.Println(coin, peer)
		if coin == 1 && peer != nil {
			fmt.Println("FLIPPED COIN sending rumor to", peer.Addr.String())
		}
	}
}

func (g *Gossiper) acceptStatusMessage(mess Messages.StatusMessage, addr *net.UDPAddr) {
	g.AddPeer(*addr)
	g.printDebugStatus(mess, *addr)
	g.peers.Decr(addr.String())
}

func (g Gossiper) printPeerList() {
	first := true
	g.peers.Lock.RLock()
	defer g.peers.Lock.RUnlock()
	for _, peer := range g.peers.Map {
		if first {
			first = false
			fmt.Print(peer)
		} else {
			fmt.Print(",", peer)
		}
	}
	fmt.Println()
}

func (g Gossiper) printDebugStatus(mess Messages.StatusMessage, addr net.UDPAddr) {
	fmt.Print("STATUS from ", addr.String())
	for _, peerStatus := range mess.Want {
		fmt.Print(" origin ", peerStatus.Identifier, " nextID ", peerStatus.NextID)
	}
	fmt.Println()
	g.printPeerList()
}

func (g Gossiper) printDebugRumor(mess Messages.RumorMessage, lastHopIP string, isFromClient bool) {
	if isFromClient {
		fmt.Println("CLIENT", mess, g.name)
	} else {
		fmt.Println("RUMOR", "origin", mess.Origin, "from", lastHopIP, "ID", mess.ID, "contents", mess.Text)
	}
	g.printPeerList()
}

func (g Gossiper) alreadySeen(id uint32, nodeName string) bool {
	g.mutex.Lock()
	_, ok := g.MessagesReceived[nodeName][id]
	g.mutex.Unlock()
	return ok
}

func (g Gossiper) storeRumorMessage(mess Messages.RumorMessage, id uint32, nodeName string) {
	g.mutex.Lock()
	if g.MessagesReceived[nodeName] == nil {
		g.MessagesReceived[nodeName] = make(map[uint32]Messages.RumorMessage)
	}
	g.MessagesReceived[nodeName][id] = mess
	g.mutex.Unlock()
}

func (g *Gossiper) antiEntropy() {
	tick := time.NewTicker(time.Minute)
	for {
		<-tick.C
		peer := g.getRandomPeer("")
		if peer != nil {
			g.sendStatusMessage(peer.Addr)
		}
	}
}

func genRouteRumor() (Messages.RumorMessage) {
	mess := Messages.RumorMessage{
		Text: "",
	}
	return mess
}

func (g *Gossiper) routeRumorDaemon() {
	tick := time.NewTicker(time.Duration(g.rtimer) * time.Second)
	for {
		<-tick.C
		g.sendRouteRumor()
	}
}

func (g *Gossiper) sendRouteRumor() {
	g.AcceptRumorMessage(genRouteRumor(), *g.gossipAddr, true)
}

func (g *Gossiper) splitComputation(mat1 gomatcore.Matrix, mat2 gomatcore.Matrix) {

}

func (g *Gossiper) splitComputations(tasks []Tasks.Task) {

}

func (g *Gossiper) acceptComputation(task Tasks.Task, addr net.Addr) bool {
	s := task.Size()
	if g.CurrentCapacity >= s {
		if _, ok := g.Pending[addr.String()]; !ok {
			g.Pending[addr.String()] = make(map[int]chan bool)
		}
		l := make(chan bool, 1)
		g.Pending[addr.String()][task.SubID] = l
		select {
		case <-l:
			//TODO faire le calcul
			g.CurrentCapacity -= s
			g.Tasks.AddTask("", task)
			return true
		case time.After(5 * time.Second):
			return false
		}
	}
	return false
}

func (g *Gossiper) listenGossiper(addr net.UDPAddr) {
	tick := time.NewTicker(time.Duration(timer * time.Second))
	IPAddress := addr.String()
	for {
		<-tick.C
		g.sendStatusMessage(addr)
		count := g.peers.Incr(IPAddress)
		if count > t1 {
			l := g.Tasks.GetTasks(addr.String())
			if l != nil && len(l) > 0 {
				g.splitComputations(l)
			}
		}
	}
}
