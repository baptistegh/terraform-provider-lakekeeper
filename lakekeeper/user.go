package lakekeeper

import (
	"net/http"
)

type (
	UserServiceInterface interface {
		GetUser(id string, options ...RequestOptionFunc) (*User, *http.Response, error)
		Whoami(options ...RequestOptionFunc) (*User, *http.Response, error)
		ProvisionUser(opts *ProvisionUserOptions, options ...RequestOptionFunc) (*User, *http.Response, error)
		DeleteUser(id string, options ...RequestOptionFunc) (*http.Response, error)
	}

	// UserService handles communication with user endpoints of the Lakekeeper API.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/user
	UserService struct {
		client *Client
	}
)

var _ UserServiceInterface = (*UserService)(nil)

// User represents a lakekeeper user
type User struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Email           *string  `json:"email,omitempty"`
	UserType        UserType `json:"user-type"`
	CreatedAt       string   `json:"created-at"`
	UpdatedAt       *string  `json:"updated-at,omitempty"`
	LastUpdatedWith string   `json:"last-updated-with"`
}

type UserType string

const (
	HumanUserType       UserType = "human"
	ApplicationUserType UserType = "application"
)

// GetUser retrieves detailed information about a specific user.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/user/operation/get_user
func (s *UserService) GetUser(id string, options ...RequestOptionFunc) (*User, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/user/"+id, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var user User

	resp, apiErr := s.client.Do(req, &user)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &user, resp, nil
}

// Whoami returns information about the user associated with the current authentication token.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/user/operation/whoami
func (s *UserService) Whoami(options ...RequestOptionFunc) (*User, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/whoami", nil, options)
	if err != nil {
		return nil, nil, err
	}

	var user User

	resp, apiErr := s.client.Do(req, &user)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &user, resp, nil
}

// ProvisionUserOptions represents ProvisionUser() options.
//
// The id must be identical to the subject in JWT tokens, prefixed with <idp-identifier>~.
// For example: oidc~1234567890 for OIDC users or kubernetes~1234567890 for Kubernetes users.
// To create users in self-service manner, do not set the id.
// The id is then extracted from the passed JWT token.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/user/operation/create_user
type ProvisionUserOptions struct {
	ID             *string   `json:"id,omitempty"`
	Email          *string   `json:"email,omitempty"`
	Name           *string   `json:"name,omitempty"`
	UpdateIfExists *bool     `json:"update-if-exists,omitempty"`
	UserType       *UserType `json:"user-type,omitempty"`
}

// ProvisionUser creates a new user or updates an existing user's metadata from the provided token.
// The token should include "profile" and "email" scopes for complete user information.
// If opts is provided, the associated user will be created
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/user/operation/create_user
func (s *UserService) ProvisionUser(opts *ProvisionUserOptions, options ...RequestOptionFunc) (*User, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "/user", opts, options)
	if err != nil {
		return nil, nil, err
	}

	var user User

	resp, apiErr := s.client.Do(req, &user)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &user, resp, nil
}

// DeleteProject permanently removes a user and all their associated permissions.
// If the user is re-registered later, their permissions will need to be re-added.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/user/operation/delete_user
func (s *UserService) DeleteUser(id string, options ...RequestOptionFunc) (*http.Response, error) {
	req, err := s.client.NewRequest(http.MethodDelete, "/user/"+id, nil, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}
