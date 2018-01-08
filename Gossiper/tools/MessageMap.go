package tools

import (
	"sync"
	"github.com/matei13/gomat/Gossiper/tools/Messages"
)

type MessageMap struct {
	Map  map[string]map[string]map[uint32]Messages.RumorMessage
	Lock *sync.RWMutex
}
