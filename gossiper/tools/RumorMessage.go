package tools

type RumorMessage struct {
	Origin      string
	PeerMessage PeerMessage
}

//NewRumorMessage -- Return an initialized RumorMessage
func NewRumorMessage(peerMessage PeerMessage) RumorMessage {
	return RumorMessage{
		PeerMessage: peerMessage,
	}
}

func (m RumorMessage) String() string {
	return m.PeerMessage.Text
}
