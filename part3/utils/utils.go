package utils

import (
	"fmt"
	"sort"
	"strconv"
)

// PeerMessage & RumorMessage
type PeerMessage struct {
	ID   uint32
	Text string
}
type RumorMessage struct {
	Origin      string
	PeerMessage PeerMessage
}

// PeerStatus & StatusPacket
type PeerStatus struct {
	Identifier string
	NextID     uint32
}
type StatusPacket struct {
	Want []PeerStatus
}

// GossiptPacket
type GossipPacket struct {
	Rumor        *RumorMessage
	StatusPacket *StatusPacket
}

// VectorClock
type vectorClockElement struct {
	Want     PeerStatus
	Messages map[uint32] /*uint32*/ *RumorMessage // PeerMessage mapped by iD
}
type VectorClock map[string]*vectorClockElement // vectorClockElement mapped by identifier

func NewVectorClock() VectorClock {
	return make(map[string]*vectorClockElement)
}

func (vc VectorClock) GetWant() []PeerStatus {
	want := make([]PeerStatus, len(vc))
	i := 0
	for _, vce := range vc {
		want[i] = vce.Want
		i++
	}
	return want
}

func foundNextID(messages map[uint32]*RumorMessage) uint32 {
	keys := make([]int, len(messages))
	i := 0
	for key := range messages {
		keys[i], _ = strconv.Atoi(fmt.Sprint(key))
		i++
	}
	sort.Ints(keys)
	for i = 0; i < len(keys); i++ {
		if keys[i] != i+1 {
			return uint32(i + 1)
		}
	}
	return uint32(len(keys) + 1)
}

func (vc VectorClock) NotifyMessage(rm *RumorMessage) {
	if _, ok := vc[rm.Origin]; !ok {
		vc[rm.Origin] = &vectorClockElement{
			Want: PeerStatus{
				Identifier: rm.Origin,
			},
			Messages: make(map[uint32]*RumorMessage),
		}
	}
	if _, ok := vc[rm.Origin].Messages[rm.PeerMessage.ID]; !ok {
		vc[rm.Origin].Messages[rm.PeerMessage.ID] = rm
		vc[rm.Origin].Want.NextID = foundNextID(vc[rm.Origin].Messages)
	}
}

// RumorMessage to send to the remote peers, bool remote peer has new messages
func (vc VectorClock) Compare(sp *StatusPacket) (*RumorMessage, bool) {
	var spVc VectorClock
	spVc = make(map[string]*vectorClockElement)
	for _, spe := range sp.Want {
		spVc[spe.Identifier] = &vectorClockElement{
			Want: spe,
		}
		if vce, ok := vc[spe.Identifier]; ok {
			if msg, ok := vce.Messages[spe.NextID]; ok {
				return msg, false
			}
			if spe.NextID > vce.Want.NextID {
				return nil, true
			}
		} else {
			vc[spe.Identifier] = &vectorClockElement{
				Want: PeerStatus{
					Identifier: spe.Identifier,
					NextID:     0,
				},
				Messages: make(map[uint32]*RumorMessage),
			}
		}
	}
	return nil, false
}

func (vc VectorClock) IsNewRumor(rm *RumorMessage) bool {
	vce, ok := vc[rm.Origin]
	if ok {
		_, ok = vce.Messages[rm.PeerMessage.ID]
	}
	return ok
}

func (vc VectorClock) GetMessages() []*RumorMessage {
	var messages []*RumorMessage
	for _, vce := range vc {
		for _, msg := range vce.Messages {
			messages = append(messages, msg)
		}
	}
	return messages
}
