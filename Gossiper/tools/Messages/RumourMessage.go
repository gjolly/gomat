package Messages

import (
	"net"
	"github.com/dedis/protobuf"
	"github.com/matei13/gomat/Daemon/gomatcore"
	"fmt"
	"github.com/matei13/gomat/matrix"
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

type MessageEncode struct {
	Origin   string
	ID       uint32
	Matrix1  []byte
	M1Col    uint32
	M1Row    uint32
	Matrix2  []byte
	M2Col    uint32
	M2Row    uint32
	M2Size   []uint
	Op       Operation
	Dest     string
	Text     string
	HopLimit uint32
}

func (m RumourMessage) String() string {
	return fmt.Sprint(m.Matrix1) + " " + fmt.Sprint(m.Matrix2)
}

// IsPrivate test if a message is private or not
func (m RumourMessage) IsPrivate() bool {
	return m.Dest != ""
}

// Send sends a Rumour message over the specified UDP connection
func (rm RumourMessage) Send(conn *net.UDPConn, addr net.UDPAddr) error {
	rmEncode, err := rm.MarshallBinary()
	if err != nil {
		return err
	}

	g := GossipMessage{Rumour: rmEncode}
	gEncode, err := protobuf.Encode(&g)
	if err != nil {
		return err
	}

	_, err = conn.WriteToUDP(gEncode, &addr)
	if err != nil {
		return err
	}
	return err
}

// MarshallBinary encodes a RumorMessage into a []byte
func (rm RumourMessage) MarshallBinary() ([]byte, error) {
	b1, err := rm.Matrix1.Mat.MarshalBinary()
	if err != nil {
		return nil, err
	}

	b2, err := rm.Matrix2.Mat.MarshalBinary()
	if err != nil {
		return nil, err
	}

	mb := MessageEncode{
		Origin:   rm.Origin,
		ID:       rm.ID,
		Matrix1:  b1,
		M1Row:    rm.Matrix1.Row,
		M1Col:    rm.Matrix1.Col,
		Matrix2:  b2,
		M2Row:    rm.Matrix2.Row,
		M2Col:    rm.Matrix2.Col,
		Op:       rm.Op,
		Dest:     rm.Dest,
		Text:     rm.Text,
		HopLimit: rm.HopLimit,
	}

	byteMessage, err := protobuf.Encode(&mb)
	if err != nil {
		return nil, err
	}

	return byteMessage, nil

}

func (me *MessageEncode) parseMatrices() (*gomatcore.SubMatrix, *gomatcore.SubMatrix, error) {
	mat1 := matrix.Matrix{}
	mat2 := matrix.Matrix{}

	err := mat1.UnmarshalBinary(me.Matrix1)
	if err != nil {
		return nil, nil, err
	}

	sm1 := gomatcore.SubMatrix{Mat: &mat1, Row: me.M1Col, Col: me.M1Row}

	err = mat2.UnmarshalBinary(me.Matrix2)
	if err != nil {
		return nil, nil, err
	}

	sm2 := gomatcore.SubMatrix{Mat: &mat2, Row: me.M2Col, Col: me.M2Row}

	return &sm1, &sm2, nil
}

// UnmarshallBinary is the opposite of MarshallBinary. It converts a []byte into a
// RumourMessage
func (rm *RumourMessage) UnmarshallBinary(buf []byte) error {
	me := &MessageEncode{}

	err := protobuf.Decode(buf, me)
	if err != nil {
		return err
	}

	sm1, sm2, err := me.parseMatrices()

	*rm = RumourMessage{
		Origin:   me.Origin,
		ID:       me.ID,
		Matrix1:  *sm1,
		Matrix2:  *sm2,
		Op:       me.Op,
		Dest:     me.Dest,
		Text:     me.Text,
		HopLimit: me.HopLimit,
	}

	return nil

}
