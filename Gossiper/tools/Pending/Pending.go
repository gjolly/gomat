package Pending

import "net"

type info struct {
	Size int
	Origin net.UDPAddr
	Chan chan bool
}

func (p *pending) getInfos() []info {
	return
}