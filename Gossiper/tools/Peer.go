package tools

import (
	"net"
)

type Peer struct {
	addr net.UDPAddr
}

func newPeer(addr net.UDPAddr) Peer {
	return Peer{
		addr: addr,
	}
}

func (p Peer) String() string {
	return p.addr.String()
}

func (p1 Peer) equals(p2 Peer) bool {
	if p1.addr.String() == p2.addr.String() {
		return true
	} else {
		return false
	}
}
