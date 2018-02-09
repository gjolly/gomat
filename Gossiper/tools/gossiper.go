package tools

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/dedis/protobuf"
	"github.com/matei13/gomat/Daemon/gomatcore"
	"github.com/matei13/gomat/Gossiper/tools/Messages"
	"github.com/matei13/gomat/Gossiper/tools/Peers"
	"github.com/matei13/gomat/Gossiper/tools/Pending"
	"github.com/matei13/gomat/Gossiper/tools/Tasks"
	"github.com/matei13/gomat/matrix"
	"log"
	"os"
)

// Gossiper -- Describe a node of a Gossip network
type Gossiper struct {
	gossipAddr      *net.UDPAddr
	UIListener      *net.UnixListener
	gossipConn      *net.UDPConn
	name            string
	peers           Peers.PeerMap
	mutex           *sync.Mutex
	PrivateMessages []Messages.RumourMessage
	MaxCapacity     int
	CurrentCapacity int
	Tasks           Tasks.TaskMap                  //Tasks[p]s: all tasks sent to p
	Pending         Pending.Pending                //Pending[k][i]: for subtask i sent from k, waiting an acknowledgement to start (true when it starts, false when it ends)
	Finished        chan bool                      //is true when the current task is finished
	TaskSize        uint32                         //number of chunks from the current task still being processed
	DoneTasks       map[uint32]gomatcore.SubMatrix //finished tasks
	FoundComputer   map[uint32]chan bool           //indicates if we have found someone to do the computations for subtask i
	Details         Details
	UnixConn        *net.UnixConn
}

type Details struct {
	rc int
	rl int
	n  int
}

const t1 = 3 //threshold for the timeouts

const timer = 30

const buffSize = 65507

// NewGossiper -- Returns a new gossiper structure
func NewGossiper(gossipPort, identifier string, peerAddrs []string, capa int) (*Gossiper, error) {
	os.Remove("/tmp/gomat.sock")
	// For UIPort
	UIAddr, err := net.ResolveUnixAddr("unix", "/tmp/gomat.sock")
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
		gossipAddr:      gossipUdpAddr,
		UIListener:      UIListener,
		gossipConn:      gossipConn,
		name:            identifier,
		peers:           Peers.PeerMap{Map: make(map[string]*Peers.Peer), Lock: &sync.RWMutex{}},
		mutex:           &sync.Mutex{},
		MaxCapacity:     capa,
		CurrentCapacity: capa,
		Tasks:           Tasks.TaskMap{Tasks: make(map[string][]Tasks.Task), Lock: &sync.RWMutex{}},
		Pending:         Pending.Pending{Infos: make(map[string]map[uint32]Pending.Info), Lock: &sync.RWMutex{}},
		Finished:        make(chan bool),
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

func (g Gossiper) sendRumourMessage(message Messages.RumourMessage, addr net.UDPAddr) error {
	message.HopLimit -= 1
	hopLimit := message.HopLimit
	rmEncode, err := message.MarshallBinary()
	gossipMessage := Messages.GossipMessage{Rumour: rmEncode}
	messEncode, err := protobuf.Encode(&gossipMessage)

	if err != nil {
		fmt.Println("error protobuf")
		return err
	}

	if addr.String() == "0.0.0.0:1234" {
		UIAddr, _ := net.ResolveUnixAddr("unix", "/tmp/gomat.sock")
		conn, _ := net.DialUnix("unix", nil, UIAddr)
		conn.Write(messEncode)
		conn.Close()
	} else if hopLimit > 0 {
		g.gossipConn.WriteToUDP(messEncode, &addr)
	}
	return nil
}

func (g Gossiper) sendStatusMessage(addr net.UDPAddr) error {
	messEncode, err := protobuf.Encode(&Messages.GossipMessage{Status: &Messages.StatusMessage{ID: -1}})
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
		unixConn, err := listener.AcceptUnix()
		if err != nil {
			log.Println("ListenUnix:", err)
		}
		g.UnixConn = unixConn
		go g.acceptUI(unixConn)
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
			addr, _ := net.ResolveUDPAddr("udp4", "0.0.0.0:1234")
			g.accept(bufferMess, addr, nbBytes, true)
		}
	}
}

// Run -- Launch the server
func (g *Gossiper) Run() {
	go g.listenUnix(g.UIListener)
	g.listenConn(g.gossipConn)
}

func (g *Gossiper) accept(buffer []byte, addr *net.UDPAddr, nbByte int, isFromClient bool) {
	mess := &Messages.GossipMessage{}
	protobuf.Decode(buffer, mess)
	if mess.Rumour != nil && len(mess.Rumour) > 0 {
		rm := &Messages.RumourMessage{}
		err := rm.UnmarshallBinary(mess.Rumour)
		if err != nil {
			log.Println(err)
			return
		}
		g.AcceptRumourMessage(*rm, *addr, isFromClient)
	} else if mess.Status != nil {
		g.acceptStatusMessage(*mess.Status, addr)
	}
}

//Callback function, call when a message is received
func (g *Gossiper) AcceptRumourMessage(mess Messages.RumourMessage, addr net.UDPAddr, isFromClient bool) {
	if isFromClient {
		g.DoneTasks = make(map[uint32]gomatcore.SubMatrix)
		g.FoundComputer = make(map[uint32]chan bool)
		mess.Origin = g.name
		// g.splitComputation(*mess.Matrix1.Mat, *mess.Matrix2.Mat, mess.Op)
		g.abcd(mess.Matrix1, mess.Matrix2, mess.Op)
		go g.merge()
	} else {
		if &mess.Matrix2 != nil {
			a := g.acceptComputation(Tasks.Task{Op: mess.Op, Mat2: mess.Matrix2, Mat1: mess.Matrix1, ID: mess.ID, Origin: addr})
			if !a {
				l := g.peers.Available(t1)
				go g.sendRumourMessage(mess, l[rand.Intn(len(l))].Addr)
			}
		} else {
			if _, ok := g.DoneTasks[mess.ID]; !ok {
				g.DoneTasks[mess.ID] = mess.Matrix1
				g.TaskSize--
				if g.TaskSize == uint32(0) {
					g.Finished <- true
				}
			}
		}
	}
}

func (g *Gossiper) acceptStatusMessage(mess Messages.StatusMessage, addr *net.UDPAddr) {
	g.AddPeer(*addr)
	if mess.ID == -1 {
		g.peers.Decr(addr.String())
	} else {
		id := uint32(mess.ID)
		s := addr.String()
		if l := g.Pending.GetChan(s, id); l != nil { //we were waiting to start the task
			l <- true
			return
		} else { // else it means that chan wants to start the task and it is not done
			select {
			case <-g.FoundComputer[uint32(mess.ID)]:
				return
			default:
				messEncode, err := protobuf.Encode(&Messages.GossipMessage{Status: &mess})
				if err != nil {
					fmt.Println("error protobuf")
				}
				g.gossipConn.WriteToUDP(messEncode, addr)
			}
		}
	}
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

func (g *Gossiper) keepSending(message Messages.RumourMessage) {
	task := Tasks.Task{
		Op:     message.Op,
		Mat2:   message.Matrix2,
		Mat1:   message.Matrix1,
		ID:     message.ID,
		Origin: *g.gossipAddr,
	}
	if g.CurrentCapacity >= task.Size() {
		g.CurrentCapacity -= task.Size()
		g.compute(task)
	} else {
		l := g.FoundComputer[message.ID]
		for {
			peerAvailable := append(g.peers.Available(t1))
			randomPeer := peerAvailable[rand.Intn(len(peerAvailable))]
			g.sendRumourMessage(message, randomPeer.Addr)
			select {
			case <-l:
				return
			default:
				time.After(5 * time.Second)
				continue
			}
		}
	}
}

func (g *Gossiper) abcd(mat1, mat2 gomatcore.SubMatrix, op Messages.Operation) {
	ansToSend := Messages.RumourMessage{}
	switch op {
	case Messages.Sum:
		ansToSend.Matrix1 = gomatcore.SubMatrix{Mat: matrix.Add(mat1.Mat, mat2.Mat), Row: mat1.Row, Col: mat2.Col}
	case Messages.Mul:
		ansToSend.Matrix1 = gomatcore.SubMatrix{Mat: matrix.Sub(mat1.Mat, mat2.Mat), Row: mat1.Row, Col: mat2.Col}
	case Messages.Sub:
		ansToSend.Matrix1 = gomatcore.SubMatrix{Mat: matrix.Mul(mat1.Mat, mat2.Mat), Row: mat1.Row, Col: mat2.Col}
	}
	rmEncode, err := ansToSend.MarshallBinary()
	if err != nil {
		log.Println("Marshall: ", rmEncode)
	}

	gm := Messages.GossipMessage{Rumour: rmEncode}
	gmEncode, err := protobuf.Encode(&gm)
	if err != nil {
		log.Println("Protobuf: ", rmEncode)
	}
	g.UnixConn.Write(gmEncode)
	g.UnixConn.Close()
}

func (g *Gossiper) splitComputation(mat1, mat2 matrix.Matrix, op Messages.Operation) {
	n := g.MaxCapacity / 2
	sMat1 := gomatcore.Split(&mat1, n)
	sMat2 := gomatcore.Split(&mat2, n)
	rl, _ := mat1.Dims()
	_, rc := mat2.Dims()
	g.Details = Details{rl: rl, rc: rc, n: n,}
	id := uint32(0)
	switch op {
	case Messages.Sum, Messages.Sub:
		for _, ssMat1 := range sMat1 {
			for _, ssMat2 := range sMat2 {
				if (ssMat1.Row == ssMat2.Row) && (ssMat1.Col == ssMat2.Col) {

					packet := Messages.RumourMessage{
						Origin:   g.name,
						ID:       id,
						Matrix1:  *ssMat1,
						Matrix2:  *ssMat2,
						Op:       op,
						HopLimit: 5,
					}
					g.FoundComputer[id] = make(chan bool, 0)
					go g.keepSending(packet)
					id++
				}
			}
		}
	case Messages.Mul:
		for _, ssMat1 := range sMat1 {
			for _, ssMat2 := range sMat2 {
				if ssMat1.Row == ssMat2.Col {
					packet := Messages.RumourMessage{
						Origin:   g.name,
						ID:       id,
						Matrix1:  *ssMat1,
						Matrix2:  *ssMat2,
						Op:       op,
						HopLimit: 5,
					}
					go g.keepSending(packet)
					g.FoundComputer[id] = make(chan bool, 0)
					id++
				}
			}
		}
	}
	g.TaskSize = id
}

func (g *Gossiper) acceptComputation(task Tasks.Task) bool {
	s := task.Size()
	if g.CurrentCapacity >= s {
		l := g.Pending.CreateChan(task)
		messEncode, err := protobuf.Encode(&Messages.GossipMessage{Status: &Messages.StatusMessage{ID: int32(task.ID)}})
		if err != nil {
			fmt.Println("error protobuf")
		}
		g.gossipConn.WriteToUDP(messEncode, &task.Origin)
		select {
		case <-l:
			l = make(chan bool, 1)
			go g.compute(task)
			g.CurrentCapacity -= s
			return true
		default:
			time.After(5 * time.Second)
			return false
		}
	}
	return false
}

func (g *Gossiper) compute(task Tasks.Task) {
	ansToSend := Messages.RumourMessage{ID: task.ID}
	switch task.Op {
	case Messages.Sum:
		ansToSend.Matrix1 = gomatcore.SubMatrix{Mat: matrix.Add(task.Mat1.Mat, task.Mat2.Mat), Row: task.Mat1.Row, Col: task.Mat2.Col}
	case Messages.Mul:
		ansToSend.Matrix1 = gomatcore.SubMatrix{Mat: matrix.Sub(task.Mat1.Mat, task.Mat2.Mat), Row: task.Mat1.Row, Col: task.Mat2.Col}
	case Messages.Sub:
		ansToSend.Matrix1 = gomatcore.SubMatrix{Mat: matrix.Mul(task.Mat1.Mat, task.Mat2.Mat), Row: task.Mat1.Row, Col: task.Mat2.Col}
	}
	g.CurrentCapacity += task.Size()
	for {
		go g.sendRumourMessage(ansToSend, task.Origin)
		ticker := time.NewTicker(5 * time.Second)
		l := g.Pending.GetChan(task.Origin.String(), task.ID)
		select {
		case <-l:
			return
		case <-ticker.C:
			continue
		}
	}
}

//listens and kills/repurposes processes associated with a dead node
func (g *Gossiper) listenGossiper(addr net.UDPAddr) {
	tick := time.NewTicker(time.Duration(timer * time.Second))
	IPAddress := addr.String()
	for {
		<-tick.C
		g.sendStatusMessage(addr)
		count := g.peers.Incr(IPAddress)
		if count > t1 {
			for _, a := range g.Pending.GetChans(addr.String()) {
				a <- true
			}
			l := g.Tasks.GetTasks(addr.String())
			if l != nil && len(l) > 0 {
				g.splitComputationList(l)
				l = make([]Tasks.Task, 0)
			}
		}
	}
}

func (g *Gossiper) splitComputationList(tasks []Tasks.Task) {
	available := g.peers.Available(t1)
	for _, t := range tasks {
		packet := Messages.RumourMessage{
			Origin:  t.Origin.String(),
			ID:      t.ID,
			Matrix1: t.Mat1,
			Matrix2: t.Mat2,
			Op:      t.Op,
		}
		randomPeer := available[rand.Intn(len(available))]
		g.sendRumourMessage(packet, randomPeer.Addr)
	}
}

// merge create the global answer from partial answers and
// sends it back to the client
func (g *Gossiper) merge() {
	<-g.Finished
	fmt.Println("After finished")
	l := make([]*gomatcore.SubMatrix, 0)
	for _, m := range g.DoneTasks {
		l = append(l, &m)
	}
	res := gomatcore.Merge(l, g.Details.rl, g.Details.rc, g.Details.n)
	mess := &Messages.RumourMessage{Matrix1: gomatcore.SubMatrix{Mat: res}}
	unixAddr, _ := net.ResolveUnixAddr("unix", "/tmp/gomat.sock")
	c, err := net.DialUnix("unix", nil, unixAddr)
	if err != nil {
		log.Println("Gossiper: Merge:", err)
	}
	defer c.Close()

	rmEncode, _ := mess.MarshallBinary()
	c.Write(rmEncode)
}
