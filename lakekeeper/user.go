package lakekeeper

import "context"

type User struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	UserType        string `json:"user-type"`
	CreatedAt       string `json:"created-at"`
	UpdatedAt       string `json:"updated-at"`
	LastUpdatedWith string `json:"last-updated-with"`
}

func (client *Client) Whoami(ctx context.Context) (*User, error) {
	var user User

	if err := client.get(ctx, "/management/v1/whoami", &user, nil); err != nil {
		return nil, err
	}

	return &user, nil
}
