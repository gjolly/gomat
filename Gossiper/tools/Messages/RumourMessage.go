package Messages

import (
	"net"
	"github.com/dedis/protobuf"
	"github.com/matei13/gomat/Daemon/gomatcore"
)

type RumourMessage struct {
	Origin   string
	ID       uint32
	Matrix1  gomatcore.SubMatrix
	Matrix2  gomatcore.SubMatrix
	Op       Operation
	Dest     string
	Text     string
	HopLimit uint32
}

func (m RumourMessage) String() string {
	return "Rumor Message"
}

func (m RumourMessage) IsPrivate() bool {
	return m.Dest != ""
}

func (rm RumourMessage) Send(conn *net.UDPConn, addr net.UDPAddr) error {
	rmEncode, err := protobuf.Encode(&rm)
	if err != nil {
		return err
	}
	_, err = conn.WriteToUDP(rmEncode, &addr)
	if err != nil {
		return err
	}
	return err
}
