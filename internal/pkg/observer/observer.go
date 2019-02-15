package observer

// Event defines an indication of a point-in-time occurrence.
type Event interface {
	EventID() string
}

// EventHandler defines a standard interface for instances that wish to list for
// the occurrence of a specific event.
type EventHandler interface {
	// Handle allows an event to be "published" to interface implementations.
	// In the "real world", error handling would likely be implemented.
	Handle(Event)
}

// Notifier is the instance being observed. Publisher is perhaps another decent
// name, but naming things is hard.
type Notifier interface {
	// Register allows an instance to register itself to listen/observe
	// events.
	Register(EventHandler)
	// Deregister allows an instance to remove itself from the collection
	// of observers/listeners.
	Deregister(EventHandler)
	// Notify publishes new events to listeners. The method is not
	// absolutely necessary, as each implementation could define this itself
	// without losing functionality.
	Notify(Event)
}

type EventNotifier struct {
	// Using a map with an empty struct allows us to keep the observers
	// unique while still keeping memory usage relatively low.
	observers map[EventHandler]struct{}
}

func NewNotifier() *EventNotifier {
	return &EventNotifier{
		observers: make(map[EventHandler]struct{}),
	}
}

func (o *EventNotifier) Register(l EventHandler) {
	o.observers[l] = struct{}{}
}

func (o *EventNotifier) Deregister(l EventHandler) {
	delete(o.observers, l)
}

func (p *EventNotifier) Notify(e Event) {
	for o := range p.observers {
		o.Handle(e)
	}
}
