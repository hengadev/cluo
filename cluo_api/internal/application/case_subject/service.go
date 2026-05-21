package caseSubjectService

import (
	"github.com/hengadev/cluo_api/internal/ports"
	"github.com/hengadev/encx"
)

type Service struct {
	repo   ports.CaseSubjectRepository
	crypto encx.CryptoService
}

func New(repo ports.CaseSubjectRepository, crypto encx.CryptoService) ports.CaseSubjectService {
	return &Service{repo: repo, crypto: crypto}
}
