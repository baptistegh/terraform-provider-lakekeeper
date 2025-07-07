package lakekeeper

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
)

// AuthSource is used to obtain access tokens.
type AuthSource interface {
	// Init is called once before making any requests.
	// If the token source needs access to client to initialize itself, it should do so here.
	Init(context.Context, *Client) error

	// Header returns an authentication header. When no error is returned, the
	// key and value should never be empty.
	Header(ctx context.Context) (key, value string, err error)
}

// OAuthTokenSource wraps an oauth2.TokenSource to implement the AuthSource interface.
type OAuthTokenSource struct {
	TokenSource oauth2.TokenSource
}

func (OAuthTokenSource) Init(context.Context, *Client) error {
	return nil
}

func (as OAuthTokenSource) Header(_ context.Context) (string, string, error) {
	t, err := as.TokenSource.Token()
	if err != nil {
		return "", "", err
	}

	return "Authorization", fmt.Sprintf("%s %s", t.TokenType, t.AccessToken), nil
}

// AccessTokenAuthSource is an AuthSource that uses a static access token.
// The token is added to the Authorization header using the Bearer scheme.
type AccessTokenAuthSource struct {
	Token string
}

func (AccessTokenAuthSource) Init(context.Context, *Client) error {
	return nil
}

func (s AccessTokenAuthSource) Header(_ context.Context) (string, string, error) {
	return "Authorization", "Bearer " + s.Token, nil
}
