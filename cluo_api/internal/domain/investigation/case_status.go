package investigation

// Status represents the status of a case
type Status string

const (
	StatusDraft      Status = "draft"
	StatusInProgress Status = "in_progress"
	StatusReady      Status = "ready"
	StatusReleased   Status = "released"
)

// IsValid checks if the Status is one of the valid constants
func (cs Status) IsValid() bool {
	switch cs {
	case StatusDraft, StatusInProgress, StatusReady, StatusReleased:
		return true
	default:
		return false
	}
}
