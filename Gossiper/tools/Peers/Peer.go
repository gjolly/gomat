package Peers

import (
	"net"
)

type Peer struct {
	Addr  net.UDPAddr
	Timer int
}

func NewPeer(addr net.UDPAddr) Peer {
	return Peer{
		Addr:  addr,
		Timer: 0,
	}
}

func (p Peer) String() string {
	return p.Addr.String()
}

func (p1 Peer) equals(p2 Peer) bool {
	if p1.Addr.String() == p2.Addr.String() {
		return true
	} else {
		return false
	}
}
