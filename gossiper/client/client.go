package main

import (
	"flag"
	"fmt"
	"net"

	"../tools"
	"github.com/dedis/protobuf"
)

func main() {
	port := flag.String("UIPort", "10000", "UIPort")
	msg := flag.String("msg", "hello", "Message")
	flag.Parse()

	udpAddr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:"+*port)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		fmt.Println(err)
		return
	}

	pmess := tools.PeerMessage{Text: *msg}
	rmess := tools.NewRumorMessage(pmess)
	mess := tools.GossipMessage{Rumor: &rmess}
	buf, err := protobuf.Encode(&mess)

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = conn.WriteToUDP(buf, udpAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
}
