package lakekeeper

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/net/publicsuffix"
)

type Config struct {
	BaseURL           string
	ClientCredentials *ClientCredentials
	Insecure          bool
	CACertFile        string
	ClientTimeout     int
	UserAgent         string
}

type Client struct {
	config       *Config
	version      *version.Version
	httpClient   *http.Client
	debug        bool
	initialLogin bool
	bootstrapped bool
}

type ClientCredentials struct {
	AuthURL      string
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

func NewClient(ctx context.Context, config *Config) (*Client, error) {
	if config.ClientCredentials.GrantType == "" {
		config.ClientCredentials.GrantType = "client_credentials"
	}

	httpClient, err := newHttpClient(config.Insecure, config.ClientTimeout, config.CACertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create http client: %v", err)
	}

	lakekeeperClient := Client{
		config:     config,
		httpClient: httpClient,
	}

	serverInfo, err := lakekeeperClient.GetServerInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting server info: %s", err)
	}

	if !lakekeeperClient.initialLogin {
		err := lakekeeperClient.login(ctx)
		if err != nil {
			return nil, fmt.Errorf("error logging in: %s", err)
		}
		lakekeeperClient.initialLogin = true
	}

	if !serverInfo.Bootstrapped {
		err := lakekeeperClient.bootstrap(ctx)
		if err != nil {
			return nil, fmt.Errorf("error bootstrapping server: %s", err)
		}
		lakekeeperClient.bootstrapped = true
	}

	if tfLog, ok := os.LookupEnv("TF_LOG"); ok {
		if tfLog == "DEBUG" {
			lakekeeperClient.debug = true
		}
	}

	return &lakekeeperClient, nil
}

func newHttpClient(tlsInsecureSkipVerify bool, clientTimeout int, caCert string) (*http.Client, error) {
	cookieJar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: tlsInsecureSkipVerify},
		Proxy:           http.ProxyFromEnvironment,
	}
	transport.MaxIdleConnsPerHost = 100

	if caCert != "" {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))
		transport.TLSClientConfig.RootCAs = caCertPool
	}

	retryClient := retryablehttp.NewClient()
	retryClient.CheckRetry = RetryPolicy
	retryClient.RetryMax = 5
	retryClient.RetryWaitMin = time.Second * 1
	retryClient.RetryWaitMax = time.Second * 60

	httpClient := retryClient.StandardClient()
	httpClient.Timeout = time.Second * time.Duration(clientTimeout)
	httpClient.Transport = transport
	httpClient.Jar = cookieJar

	return httpClient, nil
}

func (lakekeeperClient *Client) get(ctx context.Context, path string, resource interface{}, params map[string]string) error {
	body, err := lakekeeperClient.getRaw(ctx, path, params)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, resource)
}

func (lakekeeperClient *Client) getRaw(ctx context.Context, path string, params map[string]string) ([]byte, error) {
	resourceUrl := lakekeeperClient.config.BaseURL + path

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, resourceUrl, nil)
	if err != nil {
		return nil, err
	}

	if params != nil {
		query := url.Values{}
		for k, v := range params {
			query.Add(k, v)
		}
		request.URL.RawQuery = query.Encode()
	}

	body, _, err := lakekeeperClient.sendRequest(ctx, request, nil)
	return body, err
}

func (lakekeeperClient *Client) login(ctx context.Context) error {
	accessTokenData := lakekeeperClient.getAuthenticationFormData()

	tflog.Debug(ctx, "Login request", map[string]interface{}{
		"request": accessTokenData.Encode(),
	})
	accessTokenRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, lakekeeperClient.config.ClientCredentials.AuthURL, strings.NewReader(accessTokenData.Encode()))
	if err != nil {
		return err
	}

	accessTokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if lakekeeperClient.config.UserAgent != "" {
		accessTokenRequest.Header.Set("User-Agent", lakekeeperClient.config.UserAgent)
	}

	accessTokenResponse, err := lakekeeperClient.httpClient.Do(accessTokenRequest)
	if err != nil {
		return err
	}
	if accessTokenResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("error sending POST request to %s: %s", lakekeeperClient.config.ClientCredentials.AuthURL, accessTokenResponse.Status)
	}

	defer accessTokenResponse.Body.Close()

	body, _ := io.ReadAll(accessTokenResponse.Body)

	tflog.Debug(ctx, "Login response", map[string]interface{}{
		"response": string(body),
	})

	var clientCredentials ClientCredentials
	err = json.Unmarshal(body, &clientCredentials)
	if err != nil {
		return err
	}

	lakekeeperClient.config.ClientCredentials.AccessToken = clientCredentials.AccessToken
	lakekeeperClient.config.ClientCredentials.RefreshToken = clientCredentials.RefreshToken
	lakekeeperClient.config.ClientCredentials.TokenType = clientCredentials.TokenType

	info, err := lakekeeperClient.GetServerInfo(ctx)
	if err != nil {
		return err
	}

	v, err := version.NewVersion(info.Version)
	if err != nil {
		return err
	}

	lakekeeperClient.version = v

	return nil
}

func (lakekeeperClient *Client) refresh(ctx context.Context) error {
	refreshTokenData := lakekeeperClient.getAuthenticationFormData()

	tflog.Debug(ctx, "Refresh request", map[string]interface{}{
		"request": refreshTokenData.Encode(),
	})

	refreshTokenRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, lakekeeperClient.config.ClientCredentials.AuthURL, strings.NewReader(refreshTokenData.Encode()))
	if err != nil {
		return err
	}

	refreshTokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if lakekeeperClient.config.UserAgent != "" {
		refreshTokenRequest.Header.Set("User-Agent", lakekeeperClient.config.UserAgent)
	}

	refreshTokenResponse, err := lakekeeperClient.httpClient.Do(refreshTokenRequest)
	if err != nil {
		return err
	}

	defer refreshTokenResponse.Body.Close()

	body, _ := io.ReadAll(refreshTokenResponse.Body)

	tflog.Debug(ctx, "Refresh response", map[string]interface{}{
		"response": string(body),
	})

	if refreshTokenResponse.StatusCode == http.StatusBadRequest {
		tflog.Debug(ctx, "Unexpected 400, attempting to log in again")

		return lakekeeperClient.login(ctx)
	}

	var clientCredentials ClientCredentials
	err = json.Unmarshal(body, &clientCredentials)
	if err != nil {
		return err
	}

	lakekeeperClient.config.ClientCredentials.AccessToken = clientCredentials.AccessToken
	lakekeeperClient.config.ClientCredentials.RefreshToken = clientCredentials.RefreshToken
	lakekeeperClient.config.ClientCredentials.TokenType = clientCredentials.TokenType

	return nil
}

func (lakekeeperClient *Client) bootstrap(ctx context.Context) error {
	bootstrapData := `{"accept-terms-of-use":true,"is-operator":true}`

	tflog.Debug(ctx, "Bootstrap request", map[string]interface{}{
		"request": bootstrapData,
	})
	bootstrapRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, lakekeeperClient.config.BaseURL+"/management/v1/bootstrap", strings.NewReader(bootstrapData))
	if err != nil {
		return err
	}

	bootstrapRequest.Header.Set("Content-Type", "application/json")

	if lakekeeperClient.config.UserAgent != "" {
		bootstrapRequest.Header.Set("User-Agent", lakekeeperClient.config.UserAgent)
	}

	bootstrapResponse, err := lakekeeperClient.httpClient.Do(bootstrapRequest)
	if err != nil {
		return err
	}
	if bootstrapResponse.StatusCode != http.StatusNotModified {
		return fmt.Errorf("error sending POST request to %s: %s", lakekeeperClient.config.BaseURL+"/management/v1/bootstrap", bootstrapResponse.Status)
	}

	lakekeeperClient.bootstrapped = true

	return nil
}

/*
*
Sends an HTTP request and refreshes credentials on 403 or 401 errors
*/
func (lakekeeperClient *Client) sendRequest(ctx context.Context, request *http.Request, body []byte) ([]byte, string, error) {
	requestMethod := request.Method
	requestPath := request.URL.Path

	requestLogArgs := map[string]interface{}{
		"method": requestMethod,
		"path":   requestPath,
	}

	if body != nil {
		request.Body = io.NopCloser(bytes.NewReader(body))
		requestLogArgs["body"] = string(body)
	}

	tflog.Debug(ctx, "Sending request", requestLogArgs)

	lakekeeperClient.addRequestHeaders(request)

	response, err := lakekeeperClient.httpClient.Do(request)
	if err != nil {
		return nil, "", fmt.Errorf("error sending request: %v", err)
	}
	defer response.Body.Close()

	// Unauthorized: Token could have expired
	// Forbidden: After creating a realm, following GETs for the realm return 403 until you refresh
	if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden {
		tflog.Debug(ctx, "Got unexpected response, attempting refresh", map[string]interface{}{
			"status": response.Status,
		})

		err := lakekeeperClient.refresh(ctx)
		if err != nil {
			return nil, "", fmt.Errorf("error refreshing credentials: %s", err)
		}

		lakekeeperClient.addRequestHeaders(request)

		if body != nil {
			request.Body = io.NopCloser(bytes.NewReader(body))
		}
		response, err = lakekeeperClient.httpClient.Do(request)
		if err != nil {
			return nil, "", fmt.Errorf("error sending request after refresh: %v", err)
		}
		defer response.Body.Close()
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}

	responseLogArgs := map[string]interface{}{
		"status": response.Status,
	}

	if len(responseBody) != 0 && request.URL.Path != "/auth/admin/serverinfo" {
		responseLogArgs["body"] = string(responseBody)
	}

	tflog.Debug(ctx, "Received response", responseLogArgs)

	if response.StatusCode >= 400 {
		errorMessage := fmt.Sprintf("error sending %s request to %s: %s.", request.Method, request.URL.Path, response.Status)

		if len(responseBody) != 0 {
			errorMessage = fmt.Sprintf("%s Response body: %s", errorMessage, responseBody)
		}

		return nil, "", &ApiError{
			Code:    response.StatusCode,
			Message: errorMessage,
		}
	}

	return responseBody, response.Header.Get("Location"), nil
}

func (lakekeeperClient *Client) getAuthenticationFormData() url.Values {
	authenticationFormData := url.Values{}
	authenticationFormData.Set("client_id", lakekeeperClient.config.ClientCredentials.ClientID)
	authenticationFormData.Set("client_secret", lakekeeperClient.config.ClientCredentials.ClientSecret)
	authenticationFormData.Set("grant_type", lakekeeperClient.config.ClientCredentials.GrantType)

	return authenticationFormData
}

func (lakekeeperClient *Client) addRequestHeaders(request *http.Request) {
	tokenType := lakekeeperClient.config.ClientCredentials.TokenType
	accessToken := lakekeeperClient.config.ClientCredentials.AccessToken

	request.Header.Set("Authorization", fmt.Sprintf("%s %s", tokenType, accessToken))
	request.Header.Set("Accept", "application/json")

	if lakekeeperClient.config.UserAgent != "" {
		request.Header.Set("User-Agent", lakekeeperClient.config.UserAgent)
	}

	if request.Header.Get("Content-type") == "" && (request.Method == http.MethodPost || request.Method == http.MethodPut || request.Method == http.MethodDelete) {
		request.Header.Set("Content-type", "application/json")
	}
}

func RetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// do retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		return true, ctx.Err()
	}

	// 429 Too Many Requests is recoverable. Sometimes the server puts
	// a Retry-After response header to indicate when the server is
	// available to start processing request from client.
	if resp.StatusCode == http.StatusTooManyRequests {
		return true, nil
	}

	// Check the response code. We retry on 500-range responses to allow
	// the server time to recover, as 500's are typically not permanent
	// errors and may relate to outages on the server side. This will catch
	// invalid response codes as well, like 0 and 999.
	if resp.StatusCode == 0 || (resp.StatusCode >= 500 && resp.StatusCode != http.StatusNotImplemented) {
		return true, nil
	}

	return false, nil
}
