package Messages

import (
	"net"
	"github.com/dedis/protobuf"
<<<<<<< HEAD:Gossiper/tools/Messages/RumourMessage.go
	"github.com/matei13/gomat/Daemon/gomatcore"
=======
	"github.com/matei13/gomat/matrix"
	"fmt"
>>>>>>> master:Gossiper/tools/Messages/RumorMessage.go
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

<<<<<<< HEAD:Gossiper/tools/Messages/RumourMessage.go
func (m RumourMessage) String() string {
	return "Rumor Message"
=======
func (m RumorMessage) String() string {
	return fmt.Sprint(m.Matrix1) + " " + fmt.Sprint(m.Matrix2)
>>>>>>> master:Gossiper/tools/Messages/RumorMessage.go
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
