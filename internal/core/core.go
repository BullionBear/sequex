package core

type Core struct {
	timestamp int64
	orderbook map[string]Order
}

func NewCore() *Core {
	return &Core{
		timestamp: 0,
	}
}

func (c *Core) SetTimestamp(timestamp int64) {
	c.timestamp = timestamp
}

func (c *Core) GetTimestamp() int64 {
	return c.timestamp
}
