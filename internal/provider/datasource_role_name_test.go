package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestRoleExactNameLookup(t *testing.T) {
	for _, tc := range []struct {
		name      string
		pages     []string
		wantError string
	}{
		{"exact", []string{`{"roles":[{"id":"r","project-id":"p","name":"target"}]}`}, ""},
		{"missing", []string{`{"roles":[]}`}, "no role"},
		{"prefix is not exact", []string{`{"roles":[{"id":"r","project-id":"p","name":"target-extra"}]}`}, "no role"},
		{"other project", []string{`{"roles":[{"id":"r","project-id":"other","name":"target"}]}`}, "another project"},
		{"second page", []string{`{"roles":[],"next-page-token":"next"}`, `{"roles":[{"id":"r","project-id":"p","name":"target"}]}`}, ""},
		{"ambiguous across pages", []string{`{"roles":[{"id":"r","project-id":"p","name":"target"}],"next-page-token":"next"}`, `{"roles":[{"id":"s","project-id":"p","name":"target"}]}`}, "multiple roles"},
		{"repeated token", []string{`{"roles":[],"next-page-token":"next"}`, `{"roles":[],"next-page-token":"next"}`}, "repeated"},
		{"later page fails", []string{`{"roles":[{"id":"r","project-id":"p","name":"target"}],"next-page-token":"next"}`, "error"}, "list roles"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" || r.URL.Path != "/management/v1/role" || r.Header.Get("x-project-id") != "p" || r.URL.Query().Get("name") != "target" {
					t.Errorf("unexpected request %s %s", r.Method, r.URL)
				}
				if calls > 0 && r.URL.Query().Get("pageToken") != "next" {
					t.Error("missing pagination token")
				}
				if calls >= len(tc.pages) {
					t.Error("unexpected extra page")
					w.WriteHeader(500)
					return
				}
				page := tc.pages[calls]
				calls++
				w.Header().Set("Content-Type", "application/json")
				if page == "error" {
					w.WriteHeader(403)
					_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]any{"message": "denied", "type": "Forbidden", "code": 403}})
					return
				}
				_, _ = w.Write([]byte(page))
			}))
			defer server.Close()
			client, err := lakekeeper.NewClient(t.Context(), "test", server.URL, lakekeeper.WithoutRetries())
			if err != nil {
				t.Fatal(err)
			}
			role, err := findRoleByName(t.Context(), client, "p", "target")
			if tc.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantError) {
					t.Fatalf("wanted %q, got %v", tc.wantError, err)
				}
			} else if err != nil || role.ID != "r" {
				t.Fatalf("role=%v err=%v", role, err)
			}
		})
	}
}

func TestRoleReadByIDAndName(t *testing.T) {
	for _, byName := range []bool{false, true} {
		t.Run(fmt.Sprint(byName), func(t *testing.T) {
			roleJSON := `{"id":"r","project-id":"p","name":"target","created-at":"2026-10-30T10:00:00Z"}`
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if byName {
					if r.URL.Path != "/management/v1/role" {
						t.Errorf("unexpected name lookup %s", r.URL.Path)
					}
					_, _ = w.Write([]byte(`{"roles":[` + roleJSON + `]}`))
				} else {
					if r.URL.Path != "/management/v1/role/r" {
						t.Errorf("unexpected ID lookup %s", r.URL.Path)
					}
					_, _ = w.Write([]byte(roleJSON))
				}
			}))
			defer server.Close()
			client, err := lakekeeper.NewClient(t.Context(), "test", server.URL, lakekeeper.WithoutRetries())
			if err != nil {
				t.Fatal(err)
			}
			d := &LakekeeperRoleDataSource{client: client}
			var schema datasource.SchemaResponse
			d.Schema(t.Context(), datasource.SchemaRequest{}, &schema)
			state := tfsdk.State{Schema: schema.Schema}
			input := LakekeeperRoleDataSourceModel{ProjectID: types.StringValue("p")}
			if byName {
				input.Name = types.StringValue("target")
			} else {
				input.RoleID = types.StringValue("r")
			}
			if diags := state.Set(t.Context(), input); diags.HasError() {
				t.Fatal(diags)
			}
			resp := datasource.ReadResponse{State: state}
			d.Read(t.Context(), datasource.ReadRequest{Config: tfsdk.Config{Schema: schema.Schema, Raw: state.Raw}}, &resp)
			if resp.Diagnostics.HasError() {
				t.Fatal(resp.Diagnostics)
			}
			var result LakekeeperRoleDataSourceModel
			if diags := resp.State.Get(t.Context(), &result); diags.HasError() {
				t.Fatal(diags)
			}
			if result.ID.ValueString() != "p/r" || result.RoleID.ValueString() != "r" || result.Name.ValueString() != "target" {
				t.Fatalf("incorrect role state: %+v", result)
			}
		})
	}
}
