package Messages

type RumorMessage struct {
	Origin   string
	ID       uint32
	Text     string
	Dest     string
	HopLimit uint32
}

func (m RumorMessage) String() string {
	return m.Text
}

func (m RumorMessage) IsPrivate() bool {
	return m.Dest != ""
}
