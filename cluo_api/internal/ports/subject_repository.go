package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/subject"
)

type CaseSubjectRepository interface {
	ExistsCaseSubject(ctx context.Context, id uuid.UUID) (bool, error)
	CreateCaseSubject(ctx context.Context, s *subject.SubjectEncx) error
	GetCaseSubjectByID(ctx context.Context, id uuid.UUID) (*subject.SubjectEncx, error)
	ListCaseSubjects(ctx context.Context, page, pageSize int) ([]*subject.SubjectEncx, int, error)
	UpdateCaseSubject(ctx context.Context, s *subject.SubjectEncx) error
	DeleteCaseSubject(ctx context.Context, id uuid.UUID) error
}
