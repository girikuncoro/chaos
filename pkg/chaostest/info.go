package chaostest

// Info describes release information.
type Info struct {
	Description string `json:"description,omitempty"`
	Status      Status `json:"status,omitempty"`
}
