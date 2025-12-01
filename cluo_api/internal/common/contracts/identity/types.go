package identity

// Role represents a user role in the system
type Role int8

// Role enumeration values
const (
	Guest Role = iota
	Client
	Administrator
)

// Role string constants for external use
const (
	GuestStr         = "guest"
	ClientStr        = "client"
	AdministratorStr = "administrator"
)

// roleStrings maps roles to their string representations
var roleStrings = map[Role]string{
	Guest:         GuestStr,
	Client:        ClientStr,
	Administrator: AdministratorStr,
}

// stringRoles maps string representations to roles (case-insensitive)
var stringRoles = map[string]Role{
	GuestStr:         Guest,
	ClientStr:        Client,
	AdministratorStr: Administrator,
}

// roleLevels defines the hierarchy of roles (higher number = more privileges)
var roleLevels = map[Role]int{
	Guest:         0,
	Client:        1,
	Administrator: 2,
}
