package gomat

import (
	"testing"
	"os"
	"net"
	"github.com/matei13/gomat/Gossiper/tools/Messages"
	"github.com/dedis/protobuf"
	"log"
	"github.com/matei13/gomat/matrix"
)

func daemon() {
	c, err := net.Dial("unix", "/tmp/gomat.sock")
	if err != nil {
		panic(err)
	}
	requestBuf := make([]byte, 65507)

	nb, _ := c.Read(requestBuf)
	requestMess := Messages.RumorMessage{}
	protobuf.Decode(requestBuf[0:nb], &requestMess)
	log.Println(requestMess)

	r := matrix.New(2, 2, []float64{3, 3, 3, 3})
	responseMess := Messages.RumorMessage{"", 0, *r, *r, Messages.Sum, "", "", 0}
	responseBuf, _ := protobuf.Encode(&responseMess)

	c.Write(responseBuf)

}

func TestAdd(t *testing.T) {
	m1 := matrix.New(2, 2, []float64{2, 2, 2, 2})
	m2 := matrix.New(2, 2, []float64{1, 1, 1, 1})
	r, err := Add(*m1, *m2)
	if err != nil {
		panic(err)
	}
	log.Println(r)
}

func TestMain(m *testing.M) {
	go daemon()
	os.Exit(m.Run())
}
