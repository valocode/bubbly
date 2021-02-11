package events

type Event int

const (
	ResourceCreated Event = iota
	ResourceDestroyed
)

func (e Event) String() string {
	return [...]string{"CREATED", "DESTROYED"}[e]
}
