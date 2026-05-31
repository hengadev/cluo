package searchService

import (
	"context"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/sahilm/fuzzy"

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
		if r := matchCase(query, c); r != nil {
			results = append(results, *r)
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
		if m := matchField(strings.ToLower(query), c.Name, "name"); m != nil {
			results = append(results, searchDomain.Result{
				Type:    searchDomain.ResultTypeClient,
				Score:   float64(m.score),
				Item:    c,
				Matches: m.matches,
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
			if m := matchField(strings.ToLower(query), fullName, "fullName"); m != nil {
				clientName := clientNameByID[contact.ClientID]
				results = append(results, searchDomain.Result{
					Type:       searchDomain.ResultTypeContact,
					Score:      float64(m.score),
					Item:       contact,
					ClientName: &clientName,
					Matches:    m.matches,
				})
			}
		}
	}

	// sahilm/fuzzy: higher score = better match, so sort descending.
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return &searchDomain.Response{Results: results}, nil
}

type fieldMatch struct {
	score   int
	matches []searchDomain.ResultMatch
}

// matchField fuzzy-matches lowerQuery against a single field value, case-insensitively.
// Returns nil if no match.
func matchField(lowerQuery, value, key string) *fieldMatch {
	ms := fuzzy.Find(lowerQuery, []string{strings.ToLower(value)})
	if len(ms) == 0 {
		return nil
	}
	return &fieldMatch{
		score: ms[0].Score,
		matches: []searchDomain.ResultMatch{{
			Key:     key,
			Indices: collapseIndices(ms[0].MatchedIndexes),
			Value:   value,
		}},
	}
}

// matchCase returns a Result for a case, picking the best-matching field.
// Field priority: title (no penalty) > city (+10) > externalReference (+20).
func matchCase(query string, c *investigation.CaseResponse) *searchDomain.Result {
	lowerQuery := strings.ToLower(query)
	bestScore := -1
	var bestMatch []searchDomain.ResultMatch

	tryField := func(value, key string, penalty int) {
		m := matchField(lowerQuery, value, key)
		if m == nil {
			return
		}
		score := m.score - penalty
		if bestScore < 0 || score > bestScore {
			bestScore = score
			bestMatch = m.matches
		}
	}

	tryField(c.Title, "title", 0)
	if c.City != nil {
		tryField(*c.City, "city", 10)
	}
	if c.ExternalReference != nil {
		tryField(*c.ExternalReference, "externalReference", 20)
	}

	if bestScore < 0 {
		return nil
	}
	return &searchDomain.Result{
		Type:    searchDomain.ResultTypeCase,
		Score:   float64(bestScore),
		Item:    c,
		Matches: bestMatch,
	}
}

// collapseIndices merges a sorted slice of char indices into [start, end] inclusive pairs.
func collapseIndices(idxs []int) [][2]int {
	if len(idxs) == 0 {
		return nil
	}
	ranges := [][2]int{{idxs[0], idxs[0]}}
	for _, idx := range idxs[1:] {
		if idx == ranges[len(ranges)-1][1]+1 {
			ranges[len(ranges)-1][1] = idx
		} else {
			ranges = append(ranges, [2]int{idx, idx})
		}
	}
	return ranges
}
