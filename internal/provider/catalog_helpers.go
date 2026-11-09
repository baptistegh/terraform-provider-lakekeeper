package provider

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/google/uuid"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

// encodeNamespacePath validates namespace levels and joins them into the
// \x1f-separated format Lakekeeper expects. Each level must be non-empty.
func encodeNamespacePath(parts []string) (string, error) {
	for _, part := range parts {
		if part == "" {
			return "", fmt.Errorf("namespace levels must be non-empty")
		}
		if strings.ContainsAny(part, "/\x1f") {
			return "", fmt.Errorf("namespace levels cannot contain '/' or \\x1f")
		}
	}
	return strings.Join(parts, "\x1f"), nil
}

// newCatalogRequest builds an Iceberg REST catalog API request.
// It derives the catalog base URL from the client's management base URL,
// following the same pattern as go-lakekeeper's CatalogV1.
func newCatalogRequest(ctx context.Context, client *lakekeeper.Client, warehouse, suffix string, opts any, reqOpts ...core.RequestOptionFunc) (*retryablehttp.Request, error) {
	if _, err := uuid.Parse(warehouse); err != nil {
		return nil, fmt.Errorf("warehouse must be a UUID")
	}
	req, err := client.NewRequest(ctx, "GET", "/", opts, reqOpts)
	if err != nil {
		return nil, err
	}
	// Derive catalog URL: replace /management/v1 with /catalog and append /v1/{warehouse}{suffix}
	base := client.BaseURL()
	catalogBase, err := url.Parse(base.String())
	if err != nil {
		return nil, err
	}
	catalogBase.Path = strings.TrimSuffix(catalogBase.Path, "/")
	catalogBase.Path = strings.TrimSuffix(catalogBase.Path, managementv1.APIManagementVersionPath)
	catalogBase.Path = catalogBase.Path + "/catalog/v1/" + warehouse + suffix
	// Preserve the query string from the original request
	catalogBase.RawQuery = req.URL.RawQuery
	req.URL = catalogBase
	req.URL.RawPath = ""
	return req, nil
}
