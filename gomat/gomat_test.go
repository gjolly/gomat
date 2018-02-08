package gomat

import (
	"testing"
	"os"
	"net"
	"github.com/matei13/gomat/Gossiper/tools/Messages"
	"github.com/matei13/gomat/matrix"
	"time"
	"fmt"
	"github.com/matei13/gomat/Daemon/gomatcore"
	"github.com/matei13/gomat/Gossiper/tools"
	"github.com/matei13/gomat/Daemon"
)

// Daemon reads a request sent by the API and return
// a response. It uses unix socket /tmp/gomat.sock
func daemon() {
	l, err := net.Listen("unix", "/tmp/gomat.sock")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	requestBuf := make([]byte, 65507)

	c, _ := l.Accept()
	defer c.Close()

	nb, _ := c.Read(requestBuf)
	requestMess := &Messages.RumourMessage{}
	err = requestMess.UnmarshallBinary(requestBuf[0:nb])
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n+\n%v\n", requestMess.Matrix1, requestMess.Matrix2)
	r := matrix.New(2, 2, []float64{3, 3, 3, 3})
	r2 := matrix.New(2, 2, []float64{3, 3, 3, 3})
	responseMess := Messages.RumourMessage{Matrix1: gomatcore.SubMatrix{Mat: r}, Matrix2: gomatcore.SubMatrix{Mat: r2}, Op: Messages.Sum}
	responseBuf, _ := responseMess.MarshallBinary()

	c.Write(responseBuf)
}

func gomatDaemon(){
	gossiper, err := tools.NewGossiper("localhost:5000", "peerster", make([]string, 0), 5000)
	if err != nil {
		panic(err)
	}
	daemon := Daemon.Daemon{Gossiper: gossiper}
	daemon.Run()
}

func TestAdd(t *testing.T) {
	m1 := matrix.New(2, 2, []float64{2, 2, 2, 2})
	m2 := matrix.New(2, 2, []float64{1, 1, 1, 1})
	r, err := Add(m1, m2)
	if err != nil {
		panic(err)
	}
	fmt.Println("=")
	fmt.Println(r)
}

func TestMain(m *testing.M) {
	// Uncommente next line if no gomatDaemon is running
	go gomatDaemon()
	time.Sleep(time.Second)
	os.Exit(m.Run())
}
