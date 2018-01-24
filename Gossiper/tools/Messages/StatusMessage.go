package Messages

import (
	"net"
	"github.com/dedis/protobuf"
)

type StatusMessage struct {
	ID int32
}

func (sm StatusMessage) Send(conn *net.UDPConn, addr net.UDPAddr) error {
	smEncode, err := protobuf.Encode(&sm)
	if err != nil {
		return err
	}
	_, err = conn.WriteToUDP(smEncode, &addr)
	if err != nil {
		return err
	}
	return nil
}
