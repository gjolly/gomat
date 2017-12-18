package tools

import "sync"

type PeerMap struct {
	Map  map[string]Peer
	lock *sync.RWMutex
}

func (pm PeerMap) Set(k string, v Peer){
	pm.lock.Lock()
	defer pm.lock.Unlock()
	pm.Map[k] = v
}

func (pm PeerMap) Get(k string) (Peer, bool) {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	v, ok := pm.Map[k]
	return v, ok
}