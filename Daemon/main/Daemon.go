package main

import "github.com/matei13/gomat/Gossiper/tools"

type Daemon struct {
	gossiper *tools.Gossiper
}

func (d *Daemon) Run() {
	go d.gossiper.RunServer("80808")
	d.gossiper.Run()
}