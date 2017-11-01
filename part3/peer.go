package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"./utils"

	"github.com/dedis/protobuf"
	"github.com/gorilla/mux"
)

type Gossiper struct {
	connClients *net.UDPConn
	connPeers   *net.UDPConn
	name        string
	knownPeers  []*net.UDPAddr
	vectorClock utils.VectorClock
}

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func contains(s []*net.UDPAddr, e *net.UDPAddr) bool {
	for _, a := range s {
		if a.String() == e.String() {
			return true
		}
	}
	return false
}

func newGossiper(addressClients, addressPeers, name string) *Gossiper {
	udpAddrClients, err := net.ResolveUDPAddr("udp4", addressClients)
	logError(err)

	udpConnClients, err := net.ListenUDP("udp4", udpAddrClients)
	logError(err)

	udpAddrPeers, err := net.ResolveUDPAddr("udp4", addressPeers)
	logError(err)

	udpConnPeers, err := net.ListenUDP("udp4", udpAddrPeers)
	logError(err)

	return &Gossiper{
		connClients: udpConnClients,
		connPeers:   udpConnPeers,
		name:        name,
		knownPeers:  nil,
		vectorClock: utils.NewVectorClock(),
	}
}

func readDecode(conn *net.UDPConn) (utils.GossipPacket, *net.UDPAddr, error) {
	var packet utils.GossipPacket
	buf := make([]byte, 20480)

	_, relayAddr, err := conn.ReadFromUDP(buf)
	if err != nil {
		return utils.GossipPacket{}, nil, err
	}

	err = protobuf.Decode(buf, &packet)
	err = nil // FIXME: Strangely it must ignore this error
	return packet, relayAddr, err
}

func (self *Gossiper) encodeSend(address *net.UDPAddr, packet utils.GossipPacket) error {
	packetBytes, err := protobuf.Encode(&packet)
	if err != nil {
		return err
	}

	_, err = self.connPeers.WriteToUDP(packetBytes, address)
	return err
}

func (self *Gossiper) sendAck(addr *net.UDPAddr) error {
	return self.encodeSend(addr, utils.GossipPacket{
		StatusPacket: &utils.StatusPacket{
			Want: self.vectorClock.GetWant(),
		},
	})
}

func (self *Gossiper) propagateRumor(packet utils.GossipPacket) {
	if len(self.knownPeers) > 0 {
		rndPeerAdrr := self.knownPeers[rand.Intn(len(self.knownPeers))]
		fmt.Println("MONGERING with", rndPeerAdrr)
		err := self.encodeSend(rndPeerAdrr, packet)
		logError(err)

		// Coin flip
		for rand.Intn(2) == 0 {
			rndPeerAdrr := self.knownPeers[rand.Intn(len(self.knownPeers))]
			fmt.Println("FLIPPED COIN sending rumor to", rndPeerAdrr)
			err := self.encodeSend(rndPeerAdrr, packet)
			logError(err)
		}
	}
}

func (self *Gossiper) listenPeers() {
	for {
		// Read and decode
		packet, relayAddr, err := readDecode(self.connPeers)
		logError(err)

		if !contains(self.knownPeers, relayAddr) {
			self.knownPeers = append(self.knownPeers, relayAddr)
		}

		var peersAddress []string
		for _, peer := range self.knownPeers {
			peersAddress = append(peersAddress, peer.String())
		}
		fmt.Println(strings.Join(peersAddress, ","))

		// Propagate rumor
		if packet.Rumor != nil && !self.vectorClock.IsNewRumor(packet.Rumor) {
			self.vectorClock.NotifyMessage(packet.Rumor)
			err = self.sendAck(relayAddr)
			logError(err)

			fmt.Println("RUMOR origin", packet.Rumor.Origin, "from", relayAddr, "ID", packet.Rumor.PeerMessage.ID, "contents", packet.Rumor.PeerMessage.Text)

			go self.propagateRumor(packet)
		} else if packet.StatusPacket != nil {
			var wantBuffer bytes.Buffer
			for _, vce := range self.vectorClock {
				wantBuffer.WriteString("origin " + vce.Want.Identifier + " nextID " + fmt.Sprint(vce.Want.NextID) + " ")
			}
			fmt.Println("STATUS from", relayAddr, wantBuffer.String())

			msg, hasNew := self.vectorClock.Compare(packet.StatusPacket)
			if msg != nil {
				err := self.encodeSend(relayAddr, utils.GossipPacket{Rumor: msg})
				logError(err)
			} else if hasNew {
				err := self.sendAck(relayAddr)
				logError(err)
			} else {
				fmt.Println("IN SYNC WITH", relayAddr)
			}
		}
	}
}

func (self *Gossiper) listenClient() {
	for {
		packet, _, err := readDecode(self.connClients)
		logError(err)

		fmt.Println("CLIENT", packet.Rumor.PeerMessage.Text, self.name)

		var peersAddress []string
		for _, peer := range self.knownPeers {
			peersAddress = append(peersAddress, peer.String())
		}
		fmt.Println(strings.Join(peersAddress, ","))

		packet.Rumor.Origin = self.name
		id := uint32(1)
		if vce, ok := self.vectorClock[self.name]; ok {
			id = vce.Want.NextID
		}
		packet.Rumor.PeerMessage.ID = uint32(id)

		self.vectorClock.NotifyMessage(packet.Rumor)
		go self.propagateRumor(packet)
	}
}

func (self *Gossiper) antiEntropy() {
	ticker := time.NewTicker(time.Second)
	for {
		rndPeerAdrr := self.knownPeers[rand.Intn(len(self.knownPeers))]
		err := self.sendAck(rndPeerAdrr)
		logError(err)
		<-ticker.C
	}
}

func (self *Gossiper) getMessagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	json.NewEncoder(w).Encode(self.vectorClock.GetMessages())
}

func (self *Gossiper) newMessageHandler(w http.ResponseWriter, r *http.Request) {
	var rm utils.RumorMessage
	json.NewDecoder(r.Body).Decode(&rm)

	log.Println("GUI CLIENT", rm.PeerMessage.Text)

	id := uint32(1)
	if vce, ok := self.vectorClock[self.name]; ok {
		id = vce.Want.NextID
	}
	rm.PeerMessage.ID = uint32(id)

	self.vectorClock.NotifyMessage(&rm)
	go self.propagateRumor(utils.GossipPacket{Rumor: &rm})
}

func main() {
	uiPort := flag.String("UIPort", "10000", "UIport")
	gossipPort := flag.String("gossipPort", "127.0.0.1:5000", "gossipPort")
	name := flag.String("name", "", "name")
	peers := flag.String("peers", "", "peers")
	flag.Parse()

	self := newGossiper(":"+*uiPort, *gossipPort, *name)

	self.vectorClock = utils.NewVectorClock()

	if *peers == "" {
		self.knownPeers = []*net.UDPAddr{}
	} else {
		peersList := strings.Split(*peers, "_")
		self.knownPeers = make([]*net.UDPAddr, len(peersList))
		for i, peer := range peersList {
			udpAddr, err := net.ResolveUDPAddr("udp4", peer)
			logError(err)
			self.knownPeers[i] = udpAddr
		}
	}

	go self.listenPeers()
	go self.listenClient()
	go self.antiEntropy()

	r := mux.NewRouter()
	r.HandleFunc("/getMessages", self.getMessagesHandler)
	r.HandleFunc("/newMessage", self.newMessageHandler)
	http.ListenAndServe(":8080", r)
}
