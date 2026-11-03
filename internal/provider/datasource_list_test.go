package provider

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const listTestWarehouse = "22222222-2222-4222-8222-222222222222"

func listTestClient(t *testing.T, handler http.HandlerFunc) *lakekeeper.Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	client, err := lakekeeper.NewClient(t.Context(), "test", server.URL, lakekeeper.WithoutRetries())
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// readListDataSource runs Read with config as input and decodes the resulting state.
func readListDataSource[M any](t *testing.T, d datasource.DataSource, config M) (M, string) {
	t.Helper()
	var schema datasource.SchemaResponse
	d.Schema(t.Context(), datasource.SchemaRequest{}, &schema)
	state := tfsdk.State{Schema: schema.Schema}
	if diags := state.Set(t.Context(), config); diags.HasError() {
		t.Fatal(diags)
	}
	resp := datasource.ReadResponse{State: state}
	d.Read(t.Context(), datasource.ReadRequest{Config: tfsdk.Config{Schema: schema.Schema, Raw: state.Raw}}, &resp)
	var result M
	if resp.Diagnostics.HasError() {
		return result, resp.Diagnostics[0].Detail()
	}
	if diags := resp.State.Get(t.Context(), &result); diags.HasError() {
		t.Fatal(diags)
	}
	return result, ""
}

func TestProjectsListSortedByName(t *testing.T) {
	client := listTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/management/v1/project-list" {
			t.Errorf("unexpected request %s", r.URL)
		}
		_, _ = w.Write([]byte(`{"projects":[{"project-id":"b","project-name":"zeta"},{"project-id":"a","project-name":"alpha"}]}`))
	})
	result, err := readListDataSource(t, &lakekeeperProjectsDataSource{client: client}, lakekeeperProjectsDataSourceModel{})
	if err != "" || len(result.Projects) != 2 || result.Projects[0].Name.ValueString() != "alpha" || result.Projects[1].ID.ValueString() != "b" {
		t.Fatalf("projects=%+v err=%s", result.Projects, err)
	}
}

func TestWarehousesList(t *testing.T) {
	for _, tc := range []struct {
		name, body, wantError string
	}{
		{"sorted with status", `{"warehouses":[{"id":"w2","project-id":"p","name":"zeta","status":"inactive","protected":false},{"id":"w1","project-id":"p","name":"alpha","status":"active","protected":true}]}`, ""},
		{"other project", `{"warehouses":[{"id":"w1","project-id":"other","name":"alpha","status":"active"}]}`, "from project other"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := listTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/management/v1/warehouse" || r.Header.Get("x-project-id") != "p" || r.URL.Query().Has("warehouseStatus[]") {
					t.Errorf("unexpected request %s", r.URL)
				}
				_, _ = w.Write([]byte(tc.body))
			})
			result, err := readListDataSource(t, &lakekeeperWarehousesDataSource{client: client}, lakekeeperWarehousesDataSourceModel{ProjectID: types.StringValue("p")})
			if tc.wantError != "" {
				if !strings.Contains(err, tc.wantError) {
					t.Fatalf("wanted %q, got %q", tc.wantError, err)
				}
				return
			}
			got := result.Warehouses
			if err != "" || len(got) != 2 || got[0].Name.ValueString() != "alpha" || !got[0].Active.ValueBool() || !got[0].Protected.ValueBool() || got[1].Active.ValueBool() {
				t.Fatalf("warehouses=%+v err=%s", got, err)
			}
		})
	}
}

func TestNamespacesList(t *testing.T) {
	for _, tc := range []struct {
		name, warehouse, parent string
		pages                   []string
		wantNames               []string
		wantError               string
	}{
		{"top level across pages", listTestWarehouse, "", []string{
			`{"namespaces":[["instrument"]],"namespace-uuids":["n2"],"next-page-token":"next"}`,
			`{"namespaces":[["analyzer"]],"namespace-uuids":["n1"],"next-page-token":null}`,
		}, []string{"analyzer", "instrument"}, ""},
		{"children of a nested parent", listTestWarehouse, "reference.instruments", []string{
			`{"namespaces":[["reference","instruments","xrd"]],"namespace-uuids":["n3"]}`,
		}, []string{"reference.instruments.xrd"}, ""},
		{"empty", listTestWarehouse, "", []string{`{"namespaces":[]}`}, []string{}, ""},
		{"uuid count mismatch", listTestWarehouse, "", []string{`{"namespaces":[["a"]],"namespace-uuids":[]}`}, nil, "1 namespaces"},
		{"outside the parent", listTestWarehouse, "reference", []string{`{"namespaces":[["other","x"]],"namespace-uuids":["n"]}`}, nil, "outside the requested parent"},
		{"repeated token", listTestWarehouse, "", []string{
			`{"namespaces":[],"next-page-token":"next"}`, `{"namespaces":[],"next-page-token":"next"}`,
		}, nil, "repeated"},
		{"later page fails", listTestWarehouse, "", []string{`{"namespaces":[],"next-page-token":"next"}`, "error"}, nil, "list namespaces"},
		{"warehouse is not a UUID", "../x", "", nil, nil, "UUID"},
		{"invalid parent", listTestWarehouse, "a..b", nil, nil, "invalid parent"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			client := listTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				q := r.URL.Query()
				if r.Method != http.MethodGet || r.URL.Path != "/catalog/v1/"+listTestWarehouse+"/namespaces" || q.Get("returnUuids") != "true" {
					t.Errorf("unexpected request %s %s", r.Method, r.URL)
				}
				if want := strings.ReplaceAll(tc.parent, ".", "\x1f"); q.Get("parent") != want {
					t.Errorf("parent=%q, want %q", q.Get("parent"), want)
				}
				if (calls > 0) != (q.Get("pageToken") == "next") {
					t.Errorf("call %d sent pageToken %q", calls, q.Get("pageToken"))
				}
				if calls >= len(tc.pages) {
					t.Error("unexpected extra page")
					w.WriteHeader(500)
					return
				}
				page := tc.pages[calls]
				calls++
				if page == "error" {
					w.WriteHeader(403)
					_, _ = w.Write([]byte(`{"error":{"message":"denied","type":"Forbidden","code":403}}`))
					return
				}
				_, _ = w.Write([]byte(page))
			})
			config := lakekeeperNamespacesDataSourceModel{WarehouseID: types.StringValue(tc.warehouse), Parent: types.StringNull()}
			if tc.parent != "" {
				config.Parent = types.StringValue(tc.parent)
			}
			result, err := readListDataSource(t, &lakekeeperNamespacesDataSource{client: client}, config)
			if tc.wantError != "" {
				if !strings.Contains(err, tc.wantError) {
					t.Fatalf("wanted %q, got %q", tc.wantError, err)
				}
				return
			}
			names := []string{}
			for _, ns := range result.Namespaces {
				names = append(names, ns.Name.ValueString())
				if ns.ID.ValueString() == "" || strings.Join(ns.Parts, ".") != ns.Name.ValueString() {
					t.Errorf("incomplete namespace %+v", ns)
				}
			}
			if err != "" || strings.Join(names, ",") != strings.Join(tc.wantNames, ",") {
				t.Fatalf("names=%v err=%s", names, err)
			}
		})
	}
}
