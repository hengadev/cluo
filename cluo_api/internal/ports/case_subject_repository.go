package ports

import (
	"context"

	"github.com/google/uuid"
)

type CaseSubjectRepository interface {
	ExistsCaseSubject(ctx context.Context, id uuid.UUID) (bool, error)
}
