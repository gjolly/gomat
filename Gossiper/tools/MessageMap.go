package tools

import (
	"sync"
	"github.com/matei13/gomat/Gossiper/tools/Messages"
)

type MessageMap struct {
	Map  map[string]map[string]map[uint32]Messages.RumorMessage
	lock *sync.RWMutex
}

func (pm PeerMap) Set(k string, v Peer) {
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
