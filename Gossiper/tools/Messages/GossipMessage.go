package Messages

type GossipMessage struct {
	Rumour    *RumourMessage
	Status   *StatusMessage
}

func (g GossipMessage) String() string {
	var str string
	if g.Rumour != nil {
		str = "Rumor message: " + g.Rumour.String()
	} else if g.Status != nil {
		str = "Status Message"
	}

	return str
}
