package Messages

import (
	"net"
	"github.com/dedis/protobuf"
	"github.com/matei13/gomat/matrix"
)

type RumorMessage struct {
	Origin   string
	ID       uint32
	Matrix1  matrix.Matrix
	Matrix2  matrix.Matrix
	Op       Operation
	Dest     string
	Text     string
	HopLimit uint32
}

func (m RumorMessage) String() string {
	return "Rumor Message"
}

func (m RumorMessage) IsPrivate() bool {
	return m.Dest != ""
}

func (rm RumorMessage) Send(conn *net.UDPConn, addr net.UDPAddr) error {
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
