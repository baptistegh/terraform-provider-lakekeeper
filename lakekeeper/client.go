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
	InitialBootstrap  bool
}

type Client struct {
	config           *Config
	version          *version.Version
	httpClient       *http.Client
	debug            bool
	initialLogin     bool
	bootstrapped     bool
	defaultProjectID string
}

type ClientCredentials struct {
	AuthURL      string
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type BootstrapRequest struct {
	AcceptTerms bool `json:"accept-terms-of-use"`
	IsOperator  bool `json:"is-operator"`
}

const projectIDHeader = "x-project-id"

func NewClient(ctx context.Context, config *Config) (*Client, error) {
	if config.ClientCredentials.GrantType == "" {
		config.ClientCredentials.GrantType = "client_credentials"
	}

	httpClient, err := newHttpClient(config.Insecure, config.ClientTimeout, config.CACertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create http client: %v", err)
	}

	client := Client{
		config:     config,
		httpClient: httpClient,
	}

	if !client.initialLogin {
		err := client.login(ctx)
		if err != nil {
			return nil, fmt.Errorf("error logging in: %s", err)
		}
		client.initialLogin = true
	}

	serverInfo, err := client.GetServerInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting server info: %s", err)
	}

	client.bootstrapped = serverInfo.Bootstrapped
	client.defaultProjectID = serverInfo.DefaultProjectID
	if !client.bootstrapped && config.InitialBootstrap {
		err := client.bootstrap(ctx)
		if err != nil {
			return nil, fmt.Errorf("error bootstrapping server: %s", err)
		}
		client.bootstrapped = true
	}

	if tfLog, ok := os.LookupEnv("TF_LOG"); ok {
		if tfLog == "DEBUG" {
			client.debug = true
		}
	}

	return &client, nil
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

func (client *Client) get(ctx context.Context, path string, resource any, params map[string]string) error {
	body, err := client.getRaw(ctx, path, client.defaultProjectID, params)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, resource)
}

func (client *Client) getWithProjectID(ctx context.Context, path, projectID string, resource any, params map[string]string) error {
	body, err := client.getRaw(ctx, path, projectID, params)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, resource)
}

func (client *Client) postWithProjectID(ctx context.Context, path string, projectID string, body []byte) ([]byte, error) {
	resourceUrl := client.config.BaseURL + path
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, resourceUrl, nil)
	if err != nil {
		return nil, err
	}
	resp, _, err := client.sendRequest(ctx, request, projectID, body)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (client *Client) post(ctx context.Context, path string, body []byte) ([]byte, error) {
	resourceUrl := client.config.BaseURL + path

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, resourceUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, _, err := client.sendRequest(ctx, request, client.defaultProjectID, body)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (client *Client) delete(ctx context.Context, path string) error {
	resourceUrl := client.config.BaseURL + path

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, resourceUrl, nil)
	if err != nil {
		return err
	}
	_, _, err = client.sendRequest(ctx, request, client.defaultProjectID, nil)
	if err != nil {
		return err
	}
	return err
}

func (client *Client) deleteWithProjectID(ctx context.Context, path, projectID string) error {
	resourceUrl := client.config.BaseURL + path

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, resourceUrl, nil)
	if err != nil {
		return err
	}
	_, _, err = client.sendRequest(ctx, request, projectID, nil)
	if err != nil {
		return err
	}
	return err
}

func (client *Client) getRaw(ctx context.Context, path string, projectID string, params map[string]string) ([]byte, error) {
	resourceUrl := client.config.BaseURL + path

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

	body, _, err := client.sendRequest(ctx, request, projectID, nil)
	return body, err
}

// login gets a token from IDP and stores it for next uses
func (client *Client) login(ctx context.Context) error {
	accessTokenData := client.getAuthenticationFormData()

	tflog.Debug(ctx, "Login request", map[string]any{
		"request": accessTokenData.Encode(),
	})

	accessTokenRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, client.config.ClientCredentials.AuthURL, strings.NewReader(accessTokenData.Encode()))
	if err != nil {
		return err
	}

	accessTokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if client.config.UserAgent != "" {
		accessTokenRequest.Header.Set("User-Agent", client.config.UserAgent)
	}

	accessTokenResponse, err := client.httpClient.Do(accessTokenRequest)
	if err != nil {
		return err
	}
	if accessTokenResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("error sending POST request to %s: %s", client.config.ClientCredentials.AuthURL, accessTokenResponse.Status)
	}

	defer accessTokenResponse.Body.Close()

	body, _ := io.ReadAll(accessTokenResponse.Body)

	tflog.Debug(ctx, "Login response", map[string]any{
		"response": string(body),
	})

	var clientCredentials ClientCredentials
	err = json.Unmarshal(body, &clientCredentials)
	if err != nil {
		return err
	}

	client.config.ClientCredentials.AccessToken = clientCredentials.AccessToken
	client.config.ClientCredentials.RefreshToken = clientCredentials.RefreshToken
	client.config.ClientCredentials.TokenType = clientCredentials.TokenType

	info, err := client.GetServerInfo(ctx)
	if err != nil {
		return err
	}

	v, err := version.NewVersion(info.Version)
	if err != nil {
		return err
	}

	client.version = v

	return nil
}

// refresh refreshes the client token
func (client *Client) refresh(ctx context.Context) error {
	refreshTokenData := client.getAuthenticationFormData()

	tflog.Debug(ctx, "Refresh request", map[string]any{
		"request": refreshTokenData.Encode(),
	})

	refreshTokenRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, client.config.ClientCredentials.AuthURL, strings.NewReader(refreshTokenData.Encode()))
	if err != nil {
		return err
	}

	refreshTokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if client.config.UserAgent != "" {
		refreshTokenRequest.Header.Set("User-Agent", client.config.UserAgent)
	}

	refreshTokenResponse, err := client.httpClient.Do(refreshTokenRequest)
	if err != nil {
		return err
	}

	defer refreshTokenResponse.Body.Close()

	body, _ := io.ReadAll(refreshTokenResponse.Body)

	tflog.Debug(ctx, "Refresh response", map[string]any{
		"response": string(body),
	})

	if refreshTokenResponse.StatusCode == http.StatusBadRequest {
		tflog.Debug(ctx, "Unexpected 400, attempting to log in again")

		return client.login(ctx)
	}

	var clientCredentials ClientCredentials
	err = json.Unmarshal(body, &clientCredentials)
	if err != nil {
		return err
	}

	client.config.ClientCredentials.AccessToken = clientCredentials.AccessToken
	client.config.ClientCredentials.RefreshToken = clientCredentials.RefreshToken
	client.config.ClientCredentials.TokenType = clientCredentials.TokenType

	return nil
}

func (client *Client) bootstrap(ctx context.Context) error {
	bootstrapData, err := json.Marshal(BootstrapRequest{AcceptTerms: true, IsOperator: true})
	if err != nil {
		return fmt.Errorf("error creating bootstrap request, %s", err.Error())
	}

	_, err = client.post(ctx, "/management/v1/bootstrap", bootstrapData)
	if err != nil {
		return err
	}

	client.bootstrapped = true

	return nil
}

// sendRequest sends an HTTP request and refreshes credentials on 403 or 401 errors
func (client *Client) sendRequest(ctx context.Context, request *http.Request, projectID string, body []byte) ([]byte, string, error) {
	requestMethod := request.Method
	requestPath := request.URL.Path

	requestLogArgs := map[string]any{
		"method": requestMethod,
		"path":   requestPath,
	}

	if body != nil {
		request.Body = io.NopCloser(bytes.NewReader(body))
		requestLogArgs["body"] = string(body)
	}
	tflog.Debug(ctx, "Sending request", requestLogArgs)

	client.addRequestHeaders(request)
	client.setProjectIDHeader(request, projectID)

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, "", fmt.Errorf("error sending request: %v", err)
	}
	defer response.Body.Close()

	// Unauthorized: Token could have expired
	if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden {
		tflog.Debug(ctx, "Got unexpected response, attempting refresh", map[string]any{
			"status": response.Status,
		})

		err := client.refresh(ctx)
		if err != nil {
			return nil, "", fmt.Errorf("error refreshing credentials: %s", err)
		}

		client.addRequestHeaders(request)

		if body != nil {
			request.Body = io.NopCloser(bytes.NewReader(body))
		}
		response, err = client.httpClient.Do(request)
		if err != nil {
			return nil, "", fmt.Errorf("error sending request after refresh: %v", err)
		}
		defer response.Body.Close()
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}

	responseLogArgs := map[string]any{
		"status": response.Status,
	}

	if len(responseBody) != 0 {
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

func (client *Client) getAuthenticationFormData() url.Values {
	authenticationFormData := url.Values{}
	authenticationFormData.Set("client_id", client.config.ClientCredentials.ClientID)
	authenticationFormData.Set("client_secret", client.config.ClientCredentials.ClientSecret)
	authenticationFormData.Set("grant_type", client.config.ClientCredentials.GrantType)
	authenticationFormData.Set("scope", client.config.ClientCredentials.Scope)

	return authenticationFormData
}

func (client *Client) addRequestHeaders(request *http.Request) {
	tokenType := client.config.ClientCredentials.TokenType
	accessToken := client.config.ClientCredentials.AccessToken

	request.Header.Set("Authorization", fmt.Sprintf("%s %s", tokenType, accessToken))
	request.Header.Set("Accept", "application/json")

	if client.config.UserAgent != "" {
		request.Header.Set("User-Agent", client.config.UserAgent)
	}

	if request.Header.Get("Content-type") == "" && (request.Method == http.MethodPost || request.Method == http.MethodPut || request.Method == http.MethodDelete) {
		request.Header.Set("Content-type", "application/json")
	}
}

func (client *Client) setProjectIDHeader(request *http.Request, projectID string) {
	if projectID != "" {
		request.Header.Set(projectIDHeader, projectID)
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
