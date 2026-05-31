package ports

import (
	"context"

	"github.com/hengadev/cluo_api/internal/domain/search"
)

type SearchService interface {
	Search(ctx context.Context, query string) (*search.Response, error)
}
