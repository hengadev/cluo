package clientService

import (
	"github.com/hengadev/cluo_api/internal/ports"

	"github.com/hengadev/encx"
)

type Service struct {
	repo   ports.ClientRepository
	crypto encx.CryptoService
}

// New creates a new instance of the client service.
func New(user ports.ClientRepository, crypto encx.CryptoService) ports.ClientService {
	return &Service{
		repo:   user,
		crypto: crypto,
	}
}
