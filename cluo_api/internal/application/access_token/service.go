package application

import (
	"github.com/hengadev/cluo_api/internal/ports"
)

type AccessTokenService struct {
	repo ports.TokenRepository
}

func New(repo ports.TokenRepository) ports.TokenService {
	return &AccessTokenService{
		repo: repo,
	}
}
