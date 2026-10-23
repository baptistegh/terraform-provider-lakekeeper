package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestExistingNamespaceRead(t *testing.T) {
	const project = "11111111-1111-4111-8111-111111111111"
	const warehouse = "22222222-2222-4222-8222-222222222222"
	const namespace = "33333333-3333-4333-8333-333333333333"
	const name = "operational.lims.inventory.gold"
	for _, scenario := range []string{"nested", "wrong project", "wrong warehouse", "inactive", "missing namespace", "wrong path", "missing UUID", "missing managed access", "forbidden"} {
		t.Run(scenario, func(t *testing.T) {
			catalogReads := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.Header.Get("x-project-id") != project {
					t.Errorf("unexpected method or project: %s %s", r.Method, r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				switch r.URL.Path {
				case "/management/v1/warehouse/" + warehouse:
					p, id, status := project, warehouse, "active"
					if scenario == "wrong project" {
						p = namespace
					}
					if scenario == "wrong warehouse" {
						id = namespace
					}
					if scenario == "inactive" {
						status = "inactive"
					}
					_ = json.NewEncoder(w).Encode(map[string]any{"id": id, "project-id": p, "name": "renamed-warehouse", "status": status})
				case "/catalog/v1/" + warehouse + "/namespaces/operational\x1flims\x1finventory\x1fgold":
					catalogReads++
					if scenario == "missing namespace" || scenario == "forbidden" {
						code := http.StatusNotFound
						if scenario == "forbidden" {
							code = http.StatusForbidden
						}
						w.WriteHeader(code)
						_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]any{"code": code, "message": scenario, "type": "Test"}})
						return
					}
					parts, id := strings.Split(name, "."), namespace
					if scenario == "wrong path" {
						parts[3] = "bronze"
					}
					if scenario == "missing UUID" {
						id = ""
					}
					_ = json.NewEncoder(w).Encode(map[string]any{"namespace": parts, "properties": map[string]string{"namespace_id": id}})
				case "/management/v1/permissions/namespace/" + namespace:
					body := map[string]bool{}
					if scenario != "missing managed access" {
						body["managed-access"] = true
					}
					_ = json.NewEncoder(w).Encode(body)
				default:
					t.Errorf("unexpected request: %s", r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()
			ctx := context.Background()
			client, err := lakekeeper.NewClient(ctx, "test", server.URL, lakekeeper.WithoutRetries())
			if err != nil {
				t.Fatal(err)
			}
			d := &lakekeeperNamespaceDataSource{client: client}
			var schema datasource.SchemaResponse
			d.Schema(ctx, datasource.SchemaRequest{}, &schema)
			state := tfsdk.State{Schema: schema.Schema}
			model := lakekeeperNamespaceDataSourceModel{
				ProjectID: types.StringValue(project), WarehouseID: types.StringValue(warehouse), Name: types.StringValue(name),
			}
			if diags := state.Set(ctx, &model); diags.HasError() {
				t.Fatal(diags)
			}
			resp := datasource.ReadResponse{State: state}
			d.Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Schema: schema.Schema, Raw: state.Raw}}, &resp)
			if resp.Diagnostics.HasError() != (scenario != "nested") {
				t.Fatal(resp.Diagnostics)
			}
			if scenario == "nested" {
				if diags := resp.State.Get(ctx, &model); diags.HasError() {
					t.Fatal(diags)
				}
				if model.NamespaceID.ValueString() != namespace || !model.ManagedAccess.ValueBool() {
					t.Fatalf("unexpected state: %+v", model)
				}
			}
			if (scenario == "wrong project" || scenario == "wrong warehouse" || scenario == "inactive") && catalogReads != 0 {
				t.Fatal("read namespace after warehouse scope validation failed")
			}
		})
	}
}
