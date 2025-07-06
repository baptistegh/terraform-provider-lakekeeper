package lakekeeper

import (
	"context"

	"github.com/google/go-querystring/query"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

// RequestOptionFunc can be passed to all API requests to customize the API request.
type RequestOptionFunc func(*retryablehttp.Request) error

// WithHeader takes a header name and value and appends it to the request headers.
func WithHeader(name, value string) RequestOptionFunc {
	return func(req *retryablehttp.Request) error {
		req.Header.Set(name, value)
		return nil
	}
}

// WithHeaders takes a map of header name/value pairs and appends them to the
// request headers.
func WithHeaders(headers map[string]string) RequestOptionFunc {
	return func(req *retryablehttp.Request) error {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		return nil
	}
}

// WithContext runs the request with the provided context
func WithContext(ctx context.Context) RequestOptionFunc {
	return func(req *retryablehttp.Request) error {
		newCtx := copyContextValues(req.Context(), ctx)

		*req = *req.WithContext(newCtx)
		return nil
	}
}

// copyContextValues copy some context key and values in old context
func copyContextValues(oldCtx context.Context, newCtx context.Context) context.Context {
	checkRetry := checkRetryFromContext(oldCtx)

	if checkRetry != nil {
		newCtx = contextWithCheckRetry(newCtx, checkRetry)
	}

	return newCtx
}

// WithProject add the correct header in order to select a project
// for the request. The default user project is used otherwise.
func WithProject(id string) RequestOptionFunc {
	return WithHeader(HeaderProjectID, id)
}

// WithQueryParams appends query parameters derived from the given opt to the request URL.
// If `opt` is nil, the function does nothing.
// Existing query parameters with the same keys will be overwritten.
// The struct must be compatible with the
// `github.com/google/go-querystring/query` package.
func WithQueryParams(opt any) RequestOptionFunc {
	return func(req *retryablehttp.Request) error {
		if opt == nil {
			return nil
		}

		existing := req.URL.Query()

		q, err := query.Values(opt)
		if err != nil {
			return err
		}

		for key, values := range q {
			for _, v := range values {
				existing.Set(key, v)
			}
		}

		req.URL.RawQuery = existing.Encode()

		return nil
	}
}
