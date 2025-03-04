package eds

type Event struct {
	ID   string
	Name EventType
	Data interface{}
}
