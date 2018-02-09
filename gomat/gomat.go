package gomat

import (
	"net"
	"log"
	"github.com/matei13/gomat/Gossiper/tools/Messages"
	"github.com/matei13/gomat/matrix"
	"github.com/matei13/gomat/Daemon/gomatcore"
	"github.com/dedis/protobuf"
)

// askForComputation sends a computation request to the daemon via /tmp/gomat.sock
func askForComputation(m1, m2 *matrix.Matrix, operation Messages.Operation) (*matrix.Matrix, error) {
	//Resolving Unix addr to unix socket
	unixAddr, err := net.ResolveUnixAddr("unix", "/tmp/gomat.sock")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Creating the connexion
	c, err := net.DialUnix("unix", nil, unixAddr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer c.Close()

	// Creating the message
	m1Row, m1Col := m1.Dims()
	m2Row, m2Col := m2.Dims()
	rm := Messages.RumourMessage{Matrix1: gomatcore.SubMatrix{Mat: m1, Col: uint32(m1Col), Row: uint32(m1Row)}, Matrix2: gomatcore.SubMatrix{Mat: m2, Col: uint32(m2Col), Row: uint32(m2Row)}, Op: operation}

	// Encoding the message
	rmEncode, err := rm.MarshallBinary()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Sending the message to the gossiper
	gm := Messages.GossipMessage{Rumour: rmEncode}
	gmEncode, err := protobuf.Encode(&gm)
	if err != nil {
		return nil, err
	}

	_, err = c.Write(gmEncode)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Waiting for a response...
	response := make([]byte, 65507)
	nb, _, err := c.ReadFromUnix(response)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Decoding the response
	err = protobuf.Decode(response[:nb], &gm)
	if err != nil {
		log.Println("Decodage response error: ", err)
		return nil, err
	}
	responseMessage := &Messages.RumourMessage{}

	err = responseMessage.UnmarshallBinary(gm.Rumour)
	if err != nil {
		log.Println("Decodage matrix response:", err)
		return nil, err
	}

	// Returning the result
	return responseMessage.Matrix1.Mat, nil
}

// Add : Addition of two matrices
func Add(m1, m2 *matrix.Matrix) (*matrix.Matrix, error) {
	return askForComputation(m1, m2, Messages.Sum)
}

// Sub : Substruction of two matrices
func Sub(m1, m2 *matrix.Matrix) (*matrix.Matrix, error) {
	return askForComputation(m1, m2, Messages.Sub)
}

// Mult : Multiplication of two matrices
func Mult(m1, m2 *matrix.Matrix) (*matrix.Matrix, error) {
	return askForComputation(m1, m2, Messages.Mul)
}

// New creates a new Matrix
func New(r, c int, data []float64) *matrix.Matrix {
	return matrix.New(r, c, data)
}
