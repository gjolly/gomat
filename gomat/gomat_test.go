package gomat

import (
	"testing"
	"os"
	"net"
	"github.com/matei13/gomat/Gossiper/tools/Messages"
	"github.com/dedis/protobuf"
	"github.com/matei13/gomat/matrix"
	"time"
	"fmt"
)

func daemon() {
	l, err := net.Listen("unix", "/tmp/gomat.sock")
	if err != nil {
		panic(err)
	}
	requestBuf := make([]byte, 65507)

	c, _ := l.Accept()
	nb, _ := c.Read(requestBuf)
	requestMess := Messages.RumorMessage{}
	err = protobuf.Decode(requestBuf[0:nb], &requestMess)
	if err != nil {
		panic(err)
	}
	fmt.Println(requestMess)

	r := matrix.New(2, 2, []float64{3, 3, 3, 3})
	r2 := matrix.New(2, 2, []float64{3, 3, 3, 3})
	responseMess := Messages.RumorMessage{"", 0, *r, *r2, Messages.Sum, "", "", 0}
	responseBuf, _ := protobuf.Encode(&responseMess)

	c.Write(responseBuf)
	c.Close()
	l.Close()
}

func TestAdd(t *testing.T) {
	m1 := matrix.New(2, 2, []float64{2, 2, 2, 2})
	m2 := matrix.New(2, 2, []float64{1, 1, 1, 1})
	r, err := Add(m1, m2)
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
}

func TestMain(m *testing.M) {
	go daemon()
	time.Sleep(time.Second)
	os.Exit(m.Run())
}