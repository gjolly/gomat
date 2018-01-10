package Pending

import (
	"net"
	"sync"

	"github.com/matei13/gomat/Gossiper/tools/Tasks"
)

type Info struct {
	Size   int         `json:"size"`
	Origin net.UDPAddr `json:"origin"`
	Chan   chan bool
}

type Pending struct {
	Infos map[string]map[uint32]Info
	Lock  *sync.RWMutex
}

func (p *Pending) GetInfos() (tab []Info) {
	tab = make([]Info, 0)
	p.Lock.Lock()
	defer p.Lock.Unlock()
	for _, i := range p.Infos {
		for _, t := range i {
			tab = append(tab, t)
		}
	}
	return
}

func (p *Pending) GetChan(s string, i uint32) chan bool {
	p.Lock.Lock()
	defer p.Lock.Unlock()
	if l, ok := p.Infos[s]; ok {
		if a, ok := l[i]; ok {
			return a.Chan
		}
	}
	return nil
}

func (p *Pending) CreateChan(t Tasks.Task) (l chan bool) {
	p.Lock.Lock()
	defer p.Lock.Unlock()
	s := t.Origin.String()
	i := t.ID
	if _, ok := p.Infos[s]; !ok {
		p.Infos[s] = make(map[uint32]Info)
	}
	l = make(chan bool, 1)
	d1 := t.Mat1.MaxDim()
	d2 := t.Mat2.MaxDim()
	if d2 > d1 {
		d1 = d2
	}
	p.Infos[s][i] = Info{Origin: t.Origin, Chan: l, Size: d1}
	return
}
