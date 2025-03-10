package message

type Message struct {
	ID        string
	CreatedAt int64
	Content   interface{}
}
