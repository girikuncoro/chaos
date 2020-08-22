package chaostest

// Status is the status of a chaos test.
type Status string

const (
	StatusUnknown        Status = "unknown"
	StatusRunning        Status = "running"
	StatusCompleted      Status = "completed"
	StatusPendingExecute Status = "pending-execute"
)

func (s Status) String() string { return string(s) }

// IsPending determines if status is a state or a transition.
func (s Status) IsPending() bool {
	return s == StatusPendingExecute
}
