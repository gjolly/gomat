package Messages

type GossipMessage struct {
	Rumour   []byte
	Status   *StatusMessage
}

func (g GossipMessage) String() string {
	var str string
	if g.Rumour != nil {
		rm := &RumourMessage{}
		rm.UnmarshallBinary(g.Rumour)
		str = "Rumor message: " + rm.String()
	} else if g.Status != nil {
		str = "Status Message"
	}

	return str
}
