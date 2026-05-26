package document

import "testing"

func TestCanTransitionTo(t *testing.T) {
	tests := []struct {
		name    string
		from    DocumentStatus
		to      DocumentStatus
		want    bool
	}{
		// Draft → allowed
		{"draft → sent", DocumentStatusDraft, DocumentStatusSent, true},
		{"draft → cancelled", DocumentStatusDraft, DocumentStatusCancelled, true},
		{"draft → signed", DocumentStatusDraft, DocumentStatusSigned, false},
		{"draft → active", DocumentStatusDraft, DocumentStatusActive, false},
		{"draft → archived", DocumentStatusDraft, DocumentStatusArchived, false},
		{"draft → rejected", DocumentStatusDraft, DocumentStatusRejected, false},
		{"draft → expired", DocumentStatusDraft, DocumentStatusExpired, false},

		// Sent → allowed
		{"sent → signed", DocumentStatusSent, DocumentStatusSigned, true},
		{"sent → rejected", DocumentStatusSent, DocumentStatusRejected, true},
		{"sent → cancelled", DocumentStatusSent, DocumentStatusCancelled, true},
		{"sent → expired", DocumentStatusSent, DocumentStatusExpired, true},
		{"sent → draft", DocumentStatusSent, DocumentStatusDraft, false},
		{"sent → active", DocumentStatusSent, DocumentStatusActive, false},
		{"sent → archived", DocumentStatusSent, DocumentStatusArchived, false},

		// Signed → allowed
		{"signed → active", DocumentStatusSigned, DocumentStatusActive, true},
		{"signed → cancelled", DocumentStatusSigned, DocumentStatusCancelled, true},
		{"signed → draft", DocumentStatusSigned, DocumentStatusDraft, false},
		{"signed → sent", DocumentStatusSigned, DocumentStatusSent, false},
		{"signed → archived", DocumentStatusSigned, DocumentStatusArchived, false},
		{"signed → rejected", DocumentStatusSigned, DocumentStatusRejected, false},

		// Active → allowed
		{"active → archived", DocumentStatusActive, DocumentStatusArchived, true},
		{"active → cancelled", DocumentStatusActive, DocumentStatusCancelled, true},
		{"active → expired", DocumentStatusActive, DocumentStatusExpired, true},
		{"active → draft", DocumentStatusActive, DocumentStatusDraft, false},
		{"active → sent", DocumentStatusActive, DocumentStatusSent, false},
		{"active → signed", DocumentStatusActive, DocumentStatusSigned, false},

		// Archived: final state
		{"archived → draft", DocumentStatusArchived, DocumentStatusDraft, false},
		{"archived → sent", DocumentStatusArchived, DocumentStatusSent, false},
		{"archived → active", DocumentStatusArchived, DocumentStatusActive, false},
		{"archived → cancelled", DocumentStatusArchived, DocumentStatusCancelled, false},

		// Cancelled: final state
		{"cancelled → draft", DocumentStatusCancelled, DocumentStatusDraft, false},
		{"cancelled → sent", DocumentStatusCancelled, DocumentStatusSent, false},
		{"cancelled → active", DocumentStatusCancelled, DocumentStatusActive, false},
		{"cancelled → archived", DocumentStatusCancelled, DocumentStatusArchived, false},

		// Rejected → allowed
		{"rejected → draft", DocumentStatusRejected, DocumentStatusDraft, true},
		{"rejected → cancelled", DocumentStatusRejected, DocumentStatusCancelled, true},
		{"rejected → sent", DocumentStatusRejected, DocumentStatusSent, false},
		{"rejected → signed", DocumentStatusRejected, DocumentStatusSigned, false},

		// Expired → allowed
		{"expired → archived", DocumentStatusExpired, DocumentStatusArchived, true},
		{"expired → draft", DocumentStatusExpired, DocumentStatusDraft, false},
		{"expired → active", DocumentStatusExpired, DocumentStatusActive, false},

		// Invalid target status
		{"draft → invalid", DocumentStatusDraft, DocumentStatus("unknown"), false},
		{"sent → invalid", DocumentStatusSent, DocumentStatus(""), false},

		// Self-transitions
		{"draft → draft", DocumentStatusDraft, DocumentStatusDraft, false},
		{"sent → sent", DocumentStatusSent, DocumentStatusSent, false},
		{"signed → signed", DocumentStatusSigned, DocumentStatusSigned, false},
		{"active → active", DocumentStatusActive, DocumentStatusActive, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.from.CanTransitionTo(tt.to)
			if got != tt.want {
				t.Errorf("CanTransitionTo(%q → %q) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}

func TestDocumentStatus_IsValid(t *testing.T) {
	tests := []struct {
		name string
		s    DocumentStatus
		want bool
	}{
		{"draft", DocumentStatusDraft, true},
		{"sent", DocumentStatusSent, true},
		{"signed", DocumentStatusSigned, true},
		{"active", DocumentStatusActive, true},
		{"archived", DocumentStatusArchived, true},
		{"cancelled", DocumentStatusCancelled, true},
		{"rejected", DocumentStatusRejected, true},
		{"expired", DocumentStatusExpired, true},
		{"empty", DocumentStatus(""), false},
		{"unknown", DocumentStatus("pending"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsValid(); got != tt.want {
				t.Errorf("IsValid(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestDocumentStatus_IsFinal(t *testing.T) {
	tests := []struct {
		name string
		s    DocumentStatus
		want bool
	}{
		{"archived is final", DocumentStatusArchived, true},
		{"cancelled is final", DocumentStatusCancelled, true},
		{"draft is not final", DocumentStatusDraft, false},
		{"sent is not final", DocumentStatusSent, false},
		{"active is not final", DocumentStatusActive, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsFinal(); got != tt.want {
				t.Errorf("IsFinal(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}
