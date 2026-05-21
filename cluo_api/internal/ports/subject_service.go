package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/subject"
)

type CaseSubjectService interface {
	CreateCaseSubject(ctx context.Context, req *subject.CreateCaseSubjectRequest) (*subject.CaseSubjectResponse, error)
	GetCaseSubjectByID(ctx context.Context, id uuid.UUID) (*subject.CaseSubjectResponse, error)
	ListCaseSubjects(ctx context.Context, page, pageSize int) ([]*subject.CaseSubjectResponse, int, error)
	UpdateCaseSubject(ctx context.Context, req *subject.UpdateCaseSubjectRequest) (*subject.CaseSubjectResponse, error)
	DeleteCaseSubject(ctx context.Context, id uuid.UUID) error
}
