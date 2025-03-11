package core

type Core struct {
	timestamp int64
}

func NewCore() *Core {
	return &Core{
		timestamp: 0,
	}
}

func (c *Core) OnTimeUpdate(timestamp int64) {
	c.timestamp = timestamp
}

func (c *Core) GetTimestamp() int64 {
	return c.timestamp
}

func (c *Core) OnOrder(order Order) {

}
