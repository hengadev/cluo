package caseDomain

// CaseStatus represents the status of a case
type CaseStatus string

const (
	CaseStatusDraft      CaseStatus = "draft"
	CaseStatusInProgress CaseStatus = "in_progress"
	CaseStatusReady      CaseStatus = "ready"
	CaseStatusReleased   CaseStatus = "released"
)
