package lakekeeper

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
)

type User struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	UserType        string `json:"user-type"`
	CreatedAt       string `json:"created-at"`
	UpdatedAt       string `json:"updated-at"`
	LastUpdatedWith string `json:"last-updated-with"`
}

type UserCreateRequest struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	UserType       string `json:"user-type"`
	UpdateIfExists bool   `json:"update-if-exists"`
}

var (
	ValidUserTypes       = []string{"human", "application"}
	ValidUserIDPPrefixes = []string{"oidc", "kubernetes"}
)

func (client *Client) Whoami(ctx context.Context) (*User, error) {
	var user User

	if err := client.get(ctx, "/management/v1/whoami", &user, nil); err != nil {
		return nil, err
	}

	return &user, nil
}

func (client *Client) GetUserByID(ctx context.Context, id string) (*User, error) {
	if id == "" {
		return nil, fmt.Errorf("user id can not be empty")
	}
	var user User

	if err := client.get(ctx, "/management/v1/user/"+id, &user, nil); err != nil {
		return nil, err
	}

	return &user, nil
}

// NewUser creates a new user, id is required because it must match the identity provider ID
func (client *Client) NewUser(ctx context.Context, id, email, name, userType string, updateIfExists bool) (*User, error) {
	if !slices.Contains(ValidUserTypes, userType) {
		return nil, fmt.Errorf("invalid user type %s, valid values are %v", userType, ValidUserTypes)
	}
	body, err := json.Marshal(UserCreateRequest{
		ID:             id,
		Email:          email,
		Name:           name,
		UserType:       userType,
		UpdateIfExists: updateIfExists,
	})
	if err != nil {
		return nil, fmt.Errorf("could not marshal create user request, %s", err.Error())
	}

	bodyResp, err := client.post(ctx, "/management/v1/user", body)
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(bodyResp, &user); err != nil {
		return nil, fmt.Errorf("could not read create user response, %s", err.Error())
	}

	return &user, nil
}

// DeleteUser deletes a user
func (client *Client) DeleteUser(ctx context.Context, id string) error {
	err := client.delete(ctx, "/management/v1/user/"+id)
	if err != nil {
		return err
	}

	return nil
}
