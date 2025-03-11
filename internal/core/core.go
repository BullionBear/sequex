package core

import "sync"

type Core struct {
	timestamp int64
	sync.RWMutex
	m map[string]int
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
	c.Lock()
	defer c.Unlock()
	c.m[order.ID] = 1
}
