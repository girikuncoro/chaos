package chaostest

// Status is the status of a chaos test.
type Status string

const (
	StatusUnknown   Status = "unknown"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusPending   Status = "pending"
)

func (s Status) String() string { return string(s) }

// IsPending determines if status is a state or a transition.
func (s Status) IsPending() bool {
	return s == StatusPending
}
