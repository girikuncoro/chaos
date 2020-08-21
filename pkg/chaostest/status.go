package chaostest

// Status is the status of a chaos test.
type Status string

const (
	StatusUnknown        Status = "unknown"
	StatusExecuted       Status = "executed"
	StatusFail           Status = "fail"
	StatusPass           Status = "pass"
	StatusPendingExecute Status = "pending-execute"
)

func (s Status) String() string { return string(s) }
