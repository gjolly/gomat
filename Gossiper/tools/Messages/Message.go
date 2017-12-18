package Messages

import "net"

type Message interface {
	Send(conn *net.UDPConn, addr net.UDPAddr) error
}
