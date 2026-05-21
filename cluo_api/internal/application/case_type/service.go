package caseTypeService

import (
	"github.com/hengadev/cluo_api/internal/ports"
)

type Service struct {
	repo ports.CaseTypeRepository
}

func New(repo ports.CaseTypeRepository) ports.CaseTypeService {
	return &Service{repo: repo}
}
