package gomat

import (
	"net"
	"log"
	"github.com/matei13/gomat/Gossiper/tools/Messages"
	"github.com/dedis/protobuf"
	"github.com/matei13/gomat/matrix"
)

func askForComputation(m1, m2 matrix.Matrix, operation Messages.Operation) (*matrix.Matrix, error) {
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

	// Creating the message
	rm := Messages.RumourMessage{"", 0, m1, m2, operation, "", "", 0}

	// Encoding the message
	rmEncode, err := protobuf.Encode(&rm)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Sending the message to the gossiper
	_, err = c.Write(rmEncode)
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
	responseMessage := Messages.RumourMessage{}
	err = protobuf.Decode(response[0:nb], &responseMessage)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Returning the result
	return &responseMessage.Matrix1, nil
}

// Add : Addition of two matrices
func Add(m1, m2 matrix.Matrix) (*matrix.Matrix, error) {
	return askForComputation(m1, m2, Messages.Sum)
}

// Sub : Substruction of two matrices
func Sub(m1, m2 matrix.Matrix) (*matrix.Matrix, error) {
	return askForComputation(m1, m2, Messages.Subs)
}

// Mult : Multiplication of two matrices
func Mult(m1, m2 *matrix.Matrix) *matrix.Matrix {
	// TODO: As above
	return nil
}
