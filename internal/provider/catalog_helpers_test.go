package provider

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
)

func TestCatalogURLBuilding(t *testing.T) {
	warehouse := "22222222-2222-4222-8222-222222222222"
	for _, tc := range []struct {
		name, basePath, wantPath string
	}{
		{"base path without trailing slash", "/management/v1", "/catalog/v1/22222222-2222-4222-8222-222222222222/namespaces"},
		{"base path with trailing slash", "/management/v1/", "/catalog/v1/22222222-2222-4222-8222-222222222222/namespaces"},
		{"empty base path", "", "/catalog/v1/22222222-2222-4222-8222-222222222222/namespaces"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(nil)
			defer server.Close()
			baseURL := server.URL + tc.basePath
			client, err := lakekeeper.NewClient(context.Background(), "test", baseURL, lakekeeper.WithoutRetries())
			if err != nil {
				t.Fatal(err)
			}
			req, err := newCatalogRequest(context.Background(), client, warehouse, "/namespaces", nil)
			if err != nil {
				t.Fatal(err)
			}
			if req.URL.Path != tc.wantPath {
				t.Errorf("path=%q, want %q", req.URL.Path, tc.wantPath)
			}
		})
	}
}

func TestCatalogURLWithProjectHeader(t *testing.T) {
	warehouse := "22222222-2222-4222-8222-222222222222"
	project := "11111111-1111-4111-8111-111111111111"
	server := httptest.NewServer(nil)
	defer server.Close()
	client, err := lakekeeper.NewClient(context.Background(), "test", server.URL+"/management/v1", lakekeeper.WithoutRetries())
	if err != nil {
		t.Fatal(err)
	}
	req, err := newCatalogRequest(context.Background(), client, warehouse, "/namespaces/a", nil, core.WithHeader("x-project-id", project))
	if err != nil {
		t.Fatal(err)
	}
	if got := req.Header.Get("x-project-id"); got != project {
		t.Errorf("x-project-id header=%q, want %q", got, project)
	}
	if !strings.HasSuffix(req.URL.Path, "/namespaces/a") {
		t.Errorf("path=%q does not end with /namespaces/a", req.URL.Path)
	}
}

func TestEncodeNamespacePath(t *testing.T) {
	for _, tc := range []struct {
		name, want, wantError string
		parts                 []string
	}{
		{"single level", "a", "", []string{"a"}},
		{"multiple levels", "a\x1fb\x1fc", "", []string{"a", "b", "c"}},
		{"level with dot", "a.b", "", []string{"a.b"}},
		{"empty level", "", "non-empty", []string{"a", ""}},
		{"level with slash", "", "cannot contain", []string{"a/b"}},
		{"level with separator", "", "cannot contain", []string{"a\x1fb"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := encodeNamespacePath(tc.parts)
			if tc.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantError) {
					t.Fatalf("wanted error %q, got %v", tc.wantError, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Errorf("encodeNamespacePath(%v)=%q, want %q", tc.parts, got, tc.want)
			}
		})
	}
}
