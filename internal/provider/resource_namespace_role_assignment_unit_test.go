package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	lakekeeper "github.com/lakekeeper/go-lakekeeper/pkg/client"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// testNamespaceAssignment is the wire format of a namespace assignment.
type testNamespaceAssignment struct {
	Type string `json:"type"`
	Role string `json:"role,omitempty"`
	User string `json:"user,omitempty"`
}

func TestNamespaceRoleAssignmentReconcile(t *testing.T) {
	for _, tc := range []struct {
		name      string
		desired   []string
		failWrite bool
	}{
		{"replace", []string{"select"}, false}, {"revoke", nil, false}, {"partial failure", []string{"select"}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			grants := []testNamespaceAssignment{{Type: "ownership", Role: "target"}, {Type: "select", Role: "other"}, {Type: "select", User: "oidc~alice"}}
			operations := []string{}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/management/v1/permissions/namespace/ns/assignments" {
					t.Errorf("wrong scope %s", r.URL.Path)
					w.WriteHeader(400)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				if r.Method == "GET" {
					_ = json.NewEncoder(w).Encode(map[string]any{"assignments": grants})
					return
				}
				var body map[string][]testNamespaceAssignment
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				for op, changes := range body {
					operations = append(operations, op)
					if op == "writes" && tc.failWrite {
						w.WriteHeader(503)
						_, _ = w.Write([]byte(`{"error":{"code":503,"message":"unavailable","type":"Test"}}`))
						return
					}
					for _, a := range changes {
						if a.Role != "target" || a.User != "" {
							t.Errorf("changed unrelated principal: %+v", a)
						}
						if op == "writes" {
							grants = append(grants, a)
						} else {
							next := []testNamespaceAssignment{}
							for _, v := range grants {
								if v != a {
									next = append(next, v)
								}
							}
							grants = next
						}
					}
				}
				w.WriteHeader(204)
			}))
			defer server.Close()
			client, err := lakekeeper.New(ctx, server.URL, "test", lakekeeper.WithoutRetries())
			if err != nil {
				t.Fatal(err)
			}
			err = setNamespaceRoleAssignments(ctx, client, "ns", "target", tc.desired)
			if (err != nil) != tc.failWrite {
				t.Fatalf("error=%v", err)
			}
			if operations[0] != "deletes" {
				t.Fatal("must revoke before granting")
			}
			current, err := namespaceRoleAssignments(ctx, client, "ns", "target")
			if err != nil {
				t.Fatal(err)
			}
			expected := tc.desired
			if tc.failWrite {
				expected = nil
			}
			if len(current) != len(expected) {
				t.Fatalf("current %v expected %v", current, expected)
			}
			if !reflect.DeepEqual(grants[:2], []testNamespaceAssignment{{Type: "select", Role: "other"}, {Type: "select", User: "oidc~alice"}}) {
				t.Fatalf("unrelated grants changed: %v", grants)
			}
			if !tc.failWrite {
				operations = nil
				if err := setNamespaceRoleAssignments(ctx, client, "ns", "target", tc.desired); err != nil {
					t.Fatal(err)
				}
				if len(operations) != 0 {
					t.Fatal("converged apply wrote again")
				}
			}
		})
	}
}
func TestNamespaceAssignmentReadFailure(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{{"forbidden", 403, `{"error":{"code":403}}`}, {"missing", 404, `{"error":{"code":404}}`}, {"malformed", 200, `{}`}} {
		t.Run(tc.name, func(t *testing.T) {
			writes := 0
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					writes++
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer s.Close()
			c, e := lakekeeper.New(context.Background(), s.URL, "test", lakekeeper.WithoutRetries())
			if e != nil {
				t.Fatal(e)
			}
			if e := setNamespaceRoleAssignments(context.Background(), c, "n", "r", nil); e == nil {
				t.Fatal("read failure accepted")
			}
			if writes != 0 {
				t.Fatal("wrote after failed read")
			}
		})
	}
}

func TestNamespaceAssignmentResourceState(t *testing.T) {
	ctx := context.Background()
	for _, fail := range []bool{false, true} {
		t.Run(fmt.Sprint("delete_failure_", fail), func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if req.Method == "GET" {
					_, _ = w.Write([]byte(`{"assignments":[{"type":"select","role":"role"}]}`))
					return
				}
				if fail {
					w.WriteHeader(503)
					_, _ = w.Write([]byte(`{"error":{"code":503}}`))
					return
				}
				w.WriteHeader(204)
			}))
			defer s.Close()
			client, err := lakekeeper.New(ctx, s.URL, "test", lakekeeper.WithoutRetries())
			if err != nil {
				t.Fatal(err)
			}
			r := &namespaceRoleAssignmentResource{client: client}
			var schemaResp resource.SchemaResponse
			r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
			state := tfsdk.State{Schema: schemaResp.Schema}
			assignments, diags := types.SetValueFrom(ctx, types.StringType, []string{"select"})
			if diags.HasError() {
				t.Fatal(diags)
			}
			model := namespaceRoleAssignmentModel{ID: types.StringValue("p/n/role"), ProjectID: types.StringValue("p"), NamespaceID: types.StringValue("n"), RoleID: types.StringValue("role"), Assignments: assignments}
			if diags := state.Set(ctx, &model); diags.HasError() {
				t.Fatal(diags)
			}
			response := resource.DeleteResponse{State: state}
			r.Delete(ctx, resource.DeleteRequest{State: state}, &response)
			if response.Diagnostics.HasError() != fail {
				t.Fatal(response.Diagnostics)
			}
			if response.State.Raw.IsNull() == fail {
				t.Fatal("failed deletion must preserve state; successful deletion must remove it")
			}
		})
	}
}
