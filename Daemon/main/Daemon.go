package main

import "../../Gossiper/tools"

type Daemon struct {
	gossiper *tools.Gossiper
}

func (d *Daemon) Run() {
	d.gossiper.Run()
}