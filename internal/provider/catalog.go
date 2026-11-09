package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	lakekeeper "github.com/lakekeeper/go-lakekeeper/pkg/client"
)

// catalogErrorBodyLimit caps how much of an error response is quoted in the error.
const catalogErrorBodyLimit = 1024

// catalogGet sends a GET request to the Iceberg REST catalog API of
// Lakekeeper and decodes the JSON response into result. path is relative to
// /catalog/v1/, such as "{warehouse_id}/namespaces/{namespace}".
//
// The request goes through the HTTP client of the management API, so it is
// authenticated and retried the same way.
func catalogGet(ctx context.Context, client *lakekeeper.Client, project, path string, result any) error {
	u := client.BaseURL()
	u.Path = strings.TrimSuffix(u.Path, "/") + "/catalog/v1/" + path
	u.RawPath = ""

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	cfg := client.GetConfig()
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-project-id", project)
	if cfg.UserAgent != "" {
		req.Header.Set("User-Agent", cfg.UserAgent)
	}

	resp, err := cfg.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, catalogErrorBodyLimit))
		return fmt.Errorf("GET %s: %s: %s", u.Path, resp.Status, strings.TrimSpace(string(body)))
	}

	return json.NewDecoder(resp.Body).Decode(result)
}
