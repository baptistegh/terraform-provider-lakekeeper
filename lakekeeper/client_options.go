package lakekeeper

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// ClientOptionFunc can be used to customize a new Lakekeeper API client.
type ClientOptionFunc func(*Client) error

// WithCustomBackoff can be used to configure a custom backoff policy.
func WithCustomBackoff(backoff retryablehttp.Backoff) ClientOptionFunc {
	return func(c *Client) error {
		c.client.Backoff = backoff
		return nil
	}
}

// WithCustomRetry can be used to configure a custom retry policy.
func WithCustomRetry(checkRetry retryablehttp.CheckRetry) ClientOptionFunc {
	return func(c *Client) error {
		c.client.CheckRetry = checkRetry
		return nil
	}
}

// WithCustomRetryMax can be used to configure a custom maximum number of retries.
func WithCustomRetryMax(retryMax int) ClientOptionFunc {
	return func(c *Client) error {
		c.client.RetryMax = retryMax
		return nil
	}
}

// WithErrorHandler can be used to configure a custom error handler.
func WithErrorHandler(handler retryablehttp.ErrorHandler) ClientOptionFunc {
	return func(c *Client) error {
		c.client.ErrorHandler = handler
		return nil
	}
}

// WithoutRetries disables the default retry logic.
func WithoutRetries() ClientOptionFunc {
	return func(c *Client) error {
		c.disableRetries = true
		return nil
	}
}

// WithCustomRetryWaitMinMax can be used to configure a custom minimum and
// maximum time to wait between retries.
func WithCustomRetryWaitMinMax(waitMin, waitMax time.Duration) ClientOptionFunc {
	return func(c *Client) error {
		c.client.RetryWaitMin = waitMin
		c.client.RetryWaitMax = waitMax
		return nil
	}
}

// WithRequestOptions can be used to configure default request options applied to every request.
func WithRequestOptions(options ...RequestOptionFunc) ClientOptionFunc {
	return func(c *Client) error {
		c.defaultRequestOptions = append(c.defaultRequestOptions, options...)
		return nil
	}
}

// WithUserAgent can be used to configure a custom user agent.
func WithUserAgent(userAgent string) ClientOptionFunc {
	return func(c *Client) error {
		c.UserAgent = userAgent
		return nil
	}
}

// WithHTTPClient can be used to configure a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOptionFunc {
	return func(c *Client) error {
		c.client.HTTPClient = httpClient
		return nil
	}
}

// WithInitialBootstrapEnabled enables automatic server
// bootstrap on client startup.
func WithInitialBootstrapEnabled() ClientOptionFunc {
	return func(c *Client) error {
		c.bootstrap = true
		return nil
	}
}
