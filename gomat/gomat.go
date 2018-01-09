package gomat

import (
	"gonum.org/v1/gonum/mat"
	"net"
	"log"
	"github.com/matei13/gomat/Gossiper/tools/Messages"
	"github.com/dedis/protobuf"
)

// Matrix represents a matrix using the conventional storage scheme.
type Matrix struct {
	*mat.Dense
}

// New creates a new Matrix
func New(r, c int, data []float64) *Matrix {
	return &Matrix{mat.NewDense(r, c, data)}
}

func askForComputation(m1, m2 Matrix, operation int) (*Matrix, error) {
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
	rm := Messages.RumorMessage{"", 0, m1, m2, operation, "", "", 0}

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
	responseMessage := Messages.RumorMessage{}
	err = protobuf.Decode(response[0:nb], &responseMessage)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Returning the result
	return &responseMessage.Matrix1, nil
}

// Add : Addition of two matrices
func Add(m1, m2 Matrix) (*Matrix, error) {
	return askForComputation(m1, m2, Messages.Sum)
}

// Sub : Substruction of two matrices
func Sub(m1, m2 Matrix) (*Matrix, error) {
	return askForComputation(m1, m2, Messages.Sub)
}

// Mult : Multiplication of two matrices
func Mult(m1, m2 *Matrix) *Matrix {
	// TODO: As above
	return nil
}
