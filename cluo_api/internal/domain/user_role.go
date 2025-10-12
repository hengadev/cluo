package domain

// UserRole represents the role of a user
type UserRole string

const (
	UserRoleAdmin        UserRole = "admin"
	UserRoleInvestigator UserRole = "investigator"
)
