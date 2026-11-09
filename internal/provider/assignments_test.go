package provider

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	lakekeeper "github.com/lakekeeper/go-lakekeeper/pkg/client"
	"github.com/lakekeeper/go-lakekeeper/pkg/permissions"
)

// assignmentRequest records one assignment update sent to the test server.
type assignmentRequest struct {
	project string
	writes  []string
	deletes []string
}

// newAssignmentServer serves the project assignment endpoints. It answers
// failStatus to updates that contain deletes, when failStatus is set.
func newAssignmentServer(t *testing.T, failStatus int) (*lakekeeper.Client, *[]assignmentRequest) {
	t.Helper()

	requests := []assignmentRequest{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/management/v1/permissions/project/assignments" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet {
			_, _ = w.Write([]byte(`{"assignments":[]}`))
			return
		}

		var body struct {
			Writes  []struct{ Type string } `json:"writes"`
			Deletes []struct{ Type string } `json:"deletes"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		req := assignmentRequest{project: r.Header.Get("x-project-id")}
		for _, a := range body.Writes {
			req.writes = append(req.writes, a.Type)
		}
		for _, a := range body.Deletes {
			req.deletes = append(req.deletes, a.Type)
		}
		requests = append(requests, req)

		if failStatus != 0 && len(req.deletes) > 0 {
			w.WriteHeader(failStatus)
			_, _ = w.Write([]byte(`{"error":{"code":503,"message":"unavailable","type":"Test"}}`))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(server.Close)

	client, err := lakekeeper.New(t.Context(), server.URL, "test", lakekeeper.WithoutRetries())
	if err != nil {
		t.Fatal(err)
	}

	return client, &requests
}

func TestAssignmentScopeApply(t *testing.T) {
	for _, tc := range []struct {
		name           string
		added, removed []string
		want           []assignmentRequest
	}{
		{"grant before revoke", []string{"modify"}, []string{"describe"}, []assignmentRequest{
			{project: "p", writes: []string{"modify"}},
			{project: "p", deletes: []string{"describe"}},
		}},
		{"grant only", []string{"modify"}, nil, []assignmentRequest{{project: "p", writes: []string{"modify"}}}},
		{"revoke only", nil, []string{"describe"}, []assignmentRequest{{project: "p", deletes: []string{"describe"}}}},
		{"no change", nil, nil, []assignmentRequest{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client, requests := newAssignmentServer(t, 0)
			scope := projectAssignmentScope(client, "p")
			if err := scope.apply(t.Context(), permissions.PrincipalRole, "r", tc.added, tc.removed); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(*requests, tc.want) {
				t.Fatalf("requests=%+v, want %+v", *requests, tc.want)
			}
		})
	}
}

func TestAssignmentScopeApplyRevokeFails(t *testing.T) {
	client, requests := newAssignmentServer(t, http.StatusServiceUnavailable)

	err := projectAssignmentScope(client, "p").apply(t.Context(), permissions.PrincipalRole, "r", []string{"modify"}, []string{"describe"})

	var apiErr *lakekeeper.APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("want a 503 APIError, got %v", err)
	}
	if len(*requests) != 2 || (*requests)[0].writes == nil {
		t.Fatalf("the grant must be sent before the failing revoke, got %+v", *requests)
	}
}

func TestProjectAssignmentScopeRejectsEmptyProject(t *testing.T) {
	client, requests := newAssignmentServer(t, 0)
	scope := projectAssignmentScope(client, "")

	if _, err := scope.rows(t.Context()); !errors.Is(err, errNoProject) {
		t.Errorf("rows: want errNoProject, got %v", err)
	}
	if err := scope.apply(t.Context(), permissions.PrincipalRole, "r", []string{"modify"}, nil); !errors.Is(err, errNoProject) {
		t.Errorf("apply: want errNoProject, got %v", err)
	}
	if len(*requests) != 0 {
		t.Errorf("no request may reach the server's default project, got %+v", *requests)
	}
}

func TestInProject(t *testing.T) {
	if _, err := inProject(""); !errors.Is(err, errNoProject) {
		t.Errorf("want errNoProject, got %v", err)
	}
	if opt, err := inProject("p"); err != nil || opt == nil {
		t.Errorf("want an option, got %v, %v", opt, err)
	}
}
