package Daemon

import "github.com/matei13/gomat/Gossiper/tools"

type Daemon struct {
	Gossiper *tools.Gossiper
}

func (d *Daemon) Run() {
	go d.Gossiper.RunServer("8080")
	d.Gossiper.Run()
}