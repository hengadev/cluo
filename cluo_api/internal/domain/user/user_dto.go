package user

import (
	"errors"
	"time"

	"github.com/hengadev/errsx"
)

// SignInRequest represents a user sign-in request
type SignInRequest struct {
	Email    string
	Password string
}

// Valid validates the sign-in request
func (r *SignInRequest) Valid() error {
	var errs errsx.Map

	if r.Email == "" {
		errs.Set("email", errors.New("email is required"))
	}

	if r.Password == "" {
		errs.Set("password", errors.New("password is required"))
	}

	return errs.AsError()
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email    string
	Password string
}

// Valid validates the registration request
func (r *RegisterRequest) Valid() error {
	var errs errsx.Map

	if r.Email == "" {
		errs.Set("email", errors.New("email is required"))
	}

	if r.Password == "" {
		errs.Set("password", errors.New("password is required"))
	}

	if len(r.Password) > 0 && len(r.Password) < 8 {
		errs.Set("password", errors.New("password must be at least 8 characters"))
	}

	return errs.AsError()
}

// CreateSessionResponse represents a successful session creation response
type CreateSessionResponse struct {
	AccessToken          string
	RefreshToken         string
	AccessTokenExpiry    time.Time
	RefreshTokenExpiry   time.Time
	User                 *CurrentUserResponse `json:"user,omitempty"`
}

// RefreshSessionResponse represents a successful token refresh response
type RefreshSessionResponse = CreateSessionResponse

// CurrentUserResponse represents the current authenticated user
type CurrentUserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Name  string `json:"name"`
}

// UpdateNameRequest represents a request to update the user's display name
type UpdateNameRequest struct {
	Name string `json:"name"`
}
