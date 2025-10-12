package domain

// ClientType represents the type of client
type ClientType string

const (
	ClientTypePerson    ClientType = "person"
	ClientTypeLawyer    ClientType = "lawyer"
	ClientTypeInsurance ClientType = "insurance"
)
