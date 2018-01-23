package Messages

import (
	"github.com/matei13/gomat/Daemon/gomatcore"
	"log"
	"github.com/matei13/gomat/matrix"
	"testing"
)

func TestRumourMessage_MarshallBinary(t *testing.T) {

	m1 := matrix.New(2, 2, []float64{2, 2, 2, 2})
	m2 := matrix.New(2, 2, []float64{1, 1, 1, 1})

	sm1 := gomatcore.SubMatrix{Mat: m1}
	sm2 := gomatcore.SubMatrix{Mat: m2}

	// Creating the message
	rm := RumourMessage{Matrix1: sm1, Matrix2: sm2, Op: Sum}

	// Encoding the message
	rmEncode, err := rm.MarshallBinary()
	if err != nil {
		panic(err)
	}

	// Test protobuf
	testM := &RumourMessage{}
	err = testM.UnmarshallBinary(rmEncode)
	if err != nil {
		log.Println(err)
	}

	log.Println("m1 at the beginning:", m2)
	log.Println("m1 at the end:", testM.Matrix2.Mat)
}
