package gomat

import (
	"net"
	"log"
	"github.com/matei13/gomat/Gossiper/tools/Messages"
	"github.com/matei13/gomat/matrix"
	"github.com/matei13/gomat/Daemon/gomatcore"
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
	rm := Messages.RumourMessage{Matrix1: gomatcore.SubMatrix{Mat: m1}, Matrix2: gomatcore.SubMatrix{Mat: m2}, Op: operation}

	// Encoding the message
	rmEncode, err := rm.MarshallBinary()
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
	log.Println("askForComputation: size response", nb, response[:nb])
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Decoding the response
	responseMessage := &Messages.RumourMessage{}
	err = responseMessage.UnmarshallBinary(response[:nb])
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("askForComputation: response message", responseMessage)
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
