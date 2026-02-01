package ai

import (
	"github.com/google/uuid"
)

// ChatContext contains case information to inject into the chat
type ChatContext struct {
	CaseID      uuid.UUID
	CaseTitle   string
	CaseType    string
	Status      string
	ClientName  string
	Location    string
	Subjects    []ChatSubject
	RecentNotes []string
}

// ChatSubject is a simplified subject representation for chat context
type ChatSubject struct {
	Role string
	Name string
}

// estimateTokenCount provides a rough estimate of token count
// Approximately 4 characters per token for English text
func estimateTokenCount(text string) int {
	if len(text) == 0 {
		return 0
	}
	// Rough estimate: ~4 chars per token
	count := len(text) / 4
	if count < 1 {
		return 1
	}
	return count
}
