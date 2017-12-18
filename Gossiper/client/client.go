package main

import (
	"flag"
	"fmt"
	"net"
	"github.com/dedis/protobuf"
	"github.com/matei13/gomat/Gossiper/tools/Messages"
)

func main() {
	port := flag.String("UIPort", "10000", "UIPort")
	msg := flag.String("msg", "hello", "Message")
	dest := flag.String("Dest", "", "Specify a destination for a private message")
	flag.Parse()

	udpAddr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:" + *port)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		fmt.Println(err)
		return
	}

	rmess := Messages.RumorMessage{Text: *msg, Dest: *dest}
	mess := Messages.GossipMessage{Rumor: &rmess}
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
