package search

type ResultType string

const (
	ResultTypeCase    ResultType = "case"
	ResultTypeClient  ResultType = "client"
	ResultTypeContact ResultType = "contact"
)

// ResultMatch holds the character-level match positions for a single field,
// matching the Fuse.js matches format the frontend expects.
type ResultMatch struct {
	Key     string   `json:"key"`
	Indices [][2]int `json:"indices"`
	Value   string   `json:"value,omitempty"`
}

// Result is a single search hit. Item holds the concrete entity
// (CaseResponse, ClientResponse, or ContactResponse) and is marshalled
// as-is so the JSON shape matches the frontend's SearchResult type.
type Result struct {
	Type       ResultType    `json:"type"`
	Score      float64       `json:"score"`
	Item       interface{}   `json:"item"`
	ClientName *string       `json:"clientName,omitempty"`
	Matches    []ResultMatch `json:"matches,omitempty"`
}

type Response struct {
	Results []Result `json:"results"`
}
