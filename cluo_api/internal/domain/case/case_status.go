package caseDomain

// CaseStatus represents the status of a case
type CaseStatus string

const (
	CaseStatusDraft      CaseStatus = "draft"
	CaseStatusInProgress CaseStatus = "in_progress"
	CaseStatusReady      CaseStatus = "ready"
	CaseStatusReleased   CaseStatus = "released"
)

// IsValid checks if the CaseStatus is one of the valid constants
func (cs CaseStatus) IsValid() bool {
	switch cs {
	case CaseStatusDraft, CaseStatusInProgress, CaseStatusReady, CaseStatusReleased:
		return true
	default:
		return false
	}
}
