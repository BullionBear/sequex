package eds

type EventBus struct {
	subscribers map[EventType][]chan Event
}

// NewEventBus initializes an event bus
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[EventType][]chan Event),
	}
}

// Subscribe allows a consumer to listen for specific events
func (eb *EventBus) Subscribe(eventName EventType) chan Event {
	ch := make(chan Event, 1)
	eb.subscribers[eventName] = append(eb.subscribers[eventName], ch)
	return ch
}

// Publish sends an event to all subscribers
func (eb *EventBus) Publish(event Event) {
	if chans, found := eb.subscribers[event.Name]; found {
		for _, ch := range chans {
			go func(c chan Event) {
				c <- event
			}(ch)
		}
	}
}
