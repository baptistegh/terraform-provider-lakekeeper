package provider

import (
	"context"
	"fmt"
	"sort"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/sdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	managementv1 "github.com/lakekeeper/go-lakekeeper/pkg/apis/management/v1"
	lakekeeper "github.com/lakekeeper/go-lakekeeper/pkg/client"
	"github.com/lakekeeper/go-lakekeeper/pkg/permissions"
)

// assignmentScope reads, grants and revokes the direct permission assignments
// of one Lakekeeper object. T is the generated assignment union of the object
// type, such as managementv1.WarehouseAssignment.
type assignmentScope[T any] struct {
	get    func(ctx context.Context) ([]T, error)
	grant  func(ctx context.Context, assignments []T) error
	revoke func(ctx context.Context, assignments []T) error
}

func serverAssignmentScope(client *lakekeeper.Client) assignmentScope[managementv1.ServerAssignment] {
	p := client.Permissions.Server
	return assignmentScope[managementv1.ServerAssignment]{
		get: func(ctx context.Context) ([]managementv1.ServerAssignment, error) {
			resp, err := p.Assignments(ctx, nil)
			if err != nil {
				return nil, err
			}
			return resp.Assignments, nil
		},
		grant:  p.Grant,
		revoke: p.Revoke,
	}
}

func projectAssignmentScope(client *lakekeeper.Client, projectID string) assignmentScope[managementv1.ProjectAssignment] {
	p := client.Permissions.Projects
	scope, scopeErr := inProject(projectID)
	return assignmentScope[managementv1.ProjectAssignment]{
		get: func(ctx context.Context) ([]managementv1.ProjectAssignment, error) {
			if scopeErr != nil {
				return nil, scopeErr
			}
			resp, err := p.Assignments(ctx, nil, scope)
			if err != nil {
				return nil, err
			}
			return resp.Assignments, nil
		},
		grant: func(ctx context.Context, assignments []managementv1.ProjectAssignment) error {
			if scopeErr != nil {
				return scopeErr
			}
			return p.Grant(ctx, assignments, scope)
		},
		revoke: func(ctx context.Context, assignments []managementv1.ProjectAssignment) error {
			if scopeErr != nil {
				return scopeErr
			}
			return p.Revoke(ctx, assignments, scope)
		},
	}
}

func warehouseAssignmentScope(client *lakekeeper.Client, warehouseID string) assignmentScope[managementv1.WarehouseAssignment] {
	p := client.Permissions.Warehouses
	return assignmentScope[managementv1.WarehouseAssignment]{
		get: func(ctx context.Context) ([]managementv1.WarehouseAssignment, error) {
			resp, err := p.Assignments(ctx, warehouseID, nil)
			if err != nil {
				return nil, err
			}
			return resp.Assignments, nil
		},
		grant: func(ctx context.Context, assignments []managementv1.WarehouseAssignment) error {
			return p.Grant(ctx, warehouseID, assignments)
		},
		revoke: func(ctx context.Context, assignments []managementv1.WarehouseAssignment) error {
			return p.Revoke(ctx, warehouseID, assignments)
		},
	}
}

func roleAssignmentScope(client *lakekeeper.Client, roleID string) assignmentScope[managementv1.RoleAssignment] {
	p := client.Permissions.Roles
	return assignmentScope[managementv1.RoleAssignment]{
		get: func(ctx context.Context) ([]managementv1.RoleAssignment, error) {
			resp, err := p.Assignments(ctx, roleID, nil)
			if err != nil {
				return nil, err
			}
			return resp.Assignments, nil
		},
		grant: func(ctx context.Context, assignments []managementv1.RoleAssignment) error {
			return p.Grant(ctx, roleID, assignments)
		},
		revoke: func(ctx context.Context, assignments []managementv1.RoleAssignment) error {
			return p.Revoke(ctx, roleID, assignments)
		},
	}
}

// namespaceAssignmentScope uses the generated client: go-lakekeeper has no
// namespace permissions façade.
func namespaceAssignmentScope(client *lakekeeper.Client, namespaceID string) assignmentScope[managementv1.NamespaceAssignment] {
	api := client.PermissionsOpenfgaAPI
	update := func(ctx context.Context, op string, req managementv1.UpdateNamespaceAssignmentsRequest) error {
		resp, err := api.UpdateNamespaceAssignmentsById(ctx, namespaceID).UpdateNamespaceAssignmentsRequest(req).Execute()
		return lakekeeper.WrapAPIError(op, err, resp)
	}
	return assignmentScope[managementv1.NamespaceAssignment]{
		get: func(ctx context.Context) ([]managementv1.NamespaceAssignment, error) {
			out, resp, err := api.GetNamespaceAssignmentsById(ctx, namespaceID).Execute()
			if err != nil {
				return nil, lakekeeper.WrapAPIError("get namespace assignments", err, resp)
			}
			return out.Assignments, nil
		},
		grant: func(ctx context.Context, assignments []managementv1.NamespaceAssignment) error {
			return update(ctx, "grant namespace assignments", managementv1.UpdateNamespaceAssignmentsRequest{Writes: assignments})
		},
		revoke: func(ctx context.Context, assignments []managementv1.NamespaceAssignment) error {
			return update(ctx, "revoke namespace assignments", managementv1.UpdateNamespaceAssignmentsRequest{Deletes: assignments})
		},
	}
}

// rows returns every assignment of the object, in the order the server lists them.
func (s assignmentScope[T]) rows(ctx context.Context) ([]permissions.AssignmentRow, error) {
	assignments, err := s.get(ctx)
	if err != nil {
		return nil, err
	}

	rows := make([]permissions.AssignmentRow, 0, len(assignments))
	for _, a := range assignments {
		row, ok := permissions.DescribeAssignment(a)
		if !ok {
			return nil, fmt.Errorf("cannot decode assignment %+v", a)
		}
		rows = append(rows, row)
	}

	return rows, nil
}

// relations returns the sorted relations the object grants directly to one principal.
func (s assignmentScope[T]) relations(ctx context.Context, kind permissions.PrincipalKind, principalID string) ([]string, error) {
	rows, err := s.rows(ctx)
	if err != nil {
		return nil, err
	}

	principalType := principalTypeOf(kind)
	relations := []string{}
	for _, row := range rows {
		if row.PrincipalType == principalType && row.PrincipalID == principalID {
			relations = append(relations, row.Relation)
		}
	}
	sort.Strings(relations)

	return relations, nil
}

// apply grants the added relations to one principal, then revokes the removed
// ones. Granting first keeps the principal's access during the change, which
// matters when the provider's own identity is the principal. Each side is one
// request, and an empty side sends nothing. A failed revoke leaves the grants
// in place; the next refresh shows the remaining difference.
func (s assignmentScope[T]) apply(ctx context.Context, kind permissions.PrincipalKind, principalID string, added, removed []string) error {
	writes, err := buildAssignments[T](added, kind, principalID)
	if err != nil {
		return err
	}

	deletes, err := buildAssignments[T](removed, kind, principalID)
	if err != nil {
		return err
	}

	if len(writes) > 0 {
		if err := s.grant(ctx, writes); err != nil {
			return err
		}
	}

	if len(deletes) > 0 {
		return s.revoke(ctx, deletes)
	}

	return nil
}

// buildAssignments returns the assignments granting relations to one principal.
func buildAssignments[T any](relations []string, kind permissions.PrincipalKind, principalID string) ([]T, error) {
	out := make([]T, 0, len(relations))
	for _, relation := range relations {
		a, err := permissions.BuildAssignment[T](relation, kind, principalID)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}

	return out, nil
}

// principalTypeOf returns the wire name of a principal kind, as found in permissions.AssignmentRow.
func principalTypeOf(kind permissions.PrincipalKind) string {
	if kind == permissions.PrincipalRole {
		return sdk.AssigneeTypeRole
	}
	return sdk.AssigneeTypeUser
}

// stringValues returns the Go strings held by values.
func stringValues(values []types.String) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		out = append(out, v.ValueString())
	}

	return out
}
