package tools

import (
	"sync"
	"github.com/matei13/gomat/Gossiper/tools/Messages"
	"github.com/matei13/gomat/Gossiper/tools/Peers"
)

type MessageMap struct {
	Map  map[string]map[string]map[uint32]Messages.RumorMessage
	lock *sync.RWMutex
}

func (pm Peers.PeerMap) Set(k string, v Peers.Peer) {
	pm.Lock.Lock()
	defer pm.Lock.Unlock()
	pm.Map[k] = v
}

func (pm Peers.PeerMap) Get(k string) (Peers.Peer, bool) {
	pm.Lock.RLock()
	defer pm.Lock.RUnlock()
	v, ok := pm.Map[k]
	return v, ok
}
