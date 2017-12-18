package Messages

type GossipMessage struct {
	Rumor    *RumorMessage
	Status   *StatusMessage
}

func (g GossipMessage) String() string {
	var str string
	if g.Rumor != nil {
		str = "Rumor message: " + g.Rumor.String()
	} else if g.Status != nil {
		str = "Status Message"
	}

	return str
}
