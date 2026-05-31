package searchService

import (
	"context"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/lithammer/fuzzysearch/fuzzy"

	"github.com/hengadev/cluo_api/internal/domain/investigation"
	searchDomain "github.com/hengadev/cluo_api/internal/domain/search"
	"github.com/hengadev/cluo_api/internal/ports"
)

type Service struct {
	caseSvc   ports.CaseService
	clientSvc ports.ClientService
}

func New(caseSvc ports.CaseService, clientSvc ports.ClientService) *Service {
	return &Service{caseSvc: caseSvc, clientSvc: clientSvc}
}

func (s *Service) Search(ctx context.Context, query string) (*searchDomain.Response, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return &searchDomain.Response{Results: []searchDomain.Result{}}, nil
	}

	var results []searchDomain.Result

	// Cases — load all then filter in memory (data is encrypted, no SQL fuzzy possible).
	casesResp, err := s.caseSvc.List(ctx, &investigation.ListCasesRequest{Page: 1, PageSize: 1000})
	if err != nil {
		return nil, err
	}
	for _, c := range casesResp.Cases {
		if score := scoreCase(query, c); score >= 0 {
			results = append(results, searchDomain.Result{
				Type:  searchDomain.ResultTypeCase,
				Score: float64(score),
				Item:  c,
			})
		}
	}

	// Clients.
	clients, err := s.clientSvc.GetAllClients(ctx)
	if err != nil {
		return nil, err
	}

	clientNameByID := make(map[string]string, len(clients))
	for _, c := range clients {
		clientNameByID[c.ID] = c.Name
		if score := fuzzy.RankMatchFold(query, c.Name); score >= 0 {
			results = append(results, searchDomain.Result{
				Type:  searchDomain.ResultTypeClient,
				Score: float64(score),
				Item:  c,
			})
		}
	}

	// Contacts — one call per client; acceptable at single-PI dataset size.
	for _, c := range clients {
		clientID, err := uuid.Parse(c.ID)
		if err != nil {
			continue
		}
		contacts, err := s.clientSvc.GetAllContactsByClientID(ctx, clientID)
		if err != nil {
			continue
		}
		for _, contact := range contacts {
			fullName := contact.Firstname + " " + contact.Lastname
			if score := fuzzy.RankMatchFold(query, fullName); score >= 0 {
				clientName := clientNameByID[contact.ClientID]
				results = append(results, searchDomain.Result{
					Type:       searchDomain.ResultTypeContact,
					Score:      float64(score),
					Item:       contact,
					ClientName: &clientName,
				})
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score < results[j].Score
	})

	return &searchDomain.Response{Results: results}, nil
}

// scoreCase returns the best weighted fuzzy score for a case across its
// searchable fields. Returns -1 if no field matches.
// Field priority: title (no penalty) > city (+10) > externalReference (+20).
func scoreCase(query string, c *investigation.CaseResponse) int {
	best := -1

	apply := func(score, penalty int) {
		if score < 0 {
			return
		}
		weighted := score + penalty
		if best < 0 || weighted < best {
			best = weighted
		}
	}

	apply(fuzzy.RankMatchFold(query, c.Title), 0)
	if c.City != nil {
		apply(fuzzy.RankMatchFold(query, *c.City), 10)
	}
	if c.ExternalReference != nil {
		apply(fuzzy.RankMatchFold(query, *c.ExternalReference), 20)
	}

	return best
}
