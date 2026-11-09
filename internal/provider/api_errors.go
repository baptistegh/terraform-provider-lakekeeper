package provider

import (
	"errors"
	"fmt"
	"strings"

	managementv1 "github.com/lakekeeper/go-lakekeeper/pkg/apis/management/v1"
)

// apiErrorBodyLimit caps how much of a server error response is quoted in a diagnostic.
const apiErrorBodyLimit = 1024

// apiError returns err with the server's response body appended, when err is
// an error of the generated client. The generated client only reports the
// HTTP status line, while Lakekeeper explains the failure in the body.
func apiError(err error) error {
	if err == nil {
		return nil
	}

	var apiErr *managementv1.GenericOpenAPIError
	if !errors.As(err, &apiErr) {
		return err
	}

	body := strings.TrimSpace(string(apiErr.Body()))
	if body == "" {
		return err
	}
	if len(body) > apiErrorBodyLimit {
		body = body[:apiErrorBodyLimit] + "...(truncated)"
	}

	return fmt.Errorf("%w: %s", err, body)
}
