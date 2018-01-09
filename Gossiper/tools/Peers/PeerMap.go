package Peers

import (
	"sync"
)

type PeerMap struct {
	Map  map[string]*Peer
	Lock *sync.RWMutex
}

func (pm PeerMap) Set(k string, v Peer) {
	pm.Lock.Lock()
	defer pm.Lock.Unlock()
	pm.Map[k] = &v
}

func (pm PeerMap) Get(k string) (*Peer, bool) {
	pm.Lock.RLock()
	defer pm.Lock.RUnlock()
	v, ok := pm.Map[k]
	return v, ok
}

func (pm PeerMap) Incr(k string) int {
	pm.Lock.RLock()
	defer pm.Lock.RUnlock()
	_, ok := pm.Map[k]
	if ok {
		pm.Map[k].Timer++
	}
	return pm.Map[k].Timer
}

func (pm PeerMap) Decr(k string) int {
	pm.Lock.RLock()
	defer pm.Lock.RUnlock()
	_, ok := pm.Map[k]
	if ok {
		pm.Map[k].Timer--
	}
	return pm.Map[k].Timer
}

func (pm PeerMap) Available(thresh int) (peerList []*Peer) {
	peerList = make([]*Peer, 0)
	pm.Lock.RLock()
	defer pm.Lock.RUnlock()
	for _, b := range pm.Map {
		if b.Timer < thresh {
			peerList = append(peerList, b)
		}
	}
	return
}
