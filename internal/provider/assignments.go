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

// assignmentScope reads and updates the direct permission assignments of one
// Lakekeeper object. T is the generated assignment union of the object type,
// such as managementv1.WarehouseAssignment.
type assignmentScope[T any] struct {
	get    func(ctx context.Context) ([]T, error)
	update func(ctx context.Context, writes, deletes []T) error
}

func serverAssignmentScope(client *lakekeeper.Client) assignmentScope[managementv1.ServerAssignment] {
	api := client.PermissionsOpenfgaAPI
	return assignmentScope[managementv1.ServerAssignment]{
		get: func(ctx context.Context) ([]managementv1.ServerAssignment, error) {
			resp, _, err := api.GetServerAssignments(ctx).Execute()
			if err != nil {
				return nil, apiError(err)
			}
			return resp.Assignments, nil
		},
		update: func(ctx context.Context, writes, deletes []managementv1.ServerAssignment) error {
			_, err := api.UpdateServerAssignments(ctx).UpdateServerAssignmentsRequest(managementv1.UpdateServerAssignmentsRequest{
				Writes:  writes,
				Deletes: deletes,
			}).Execute()
			return apiError(err)
		},
	}
}

func projectAssignmentScope(client *lakekeeper.Client, projectID string) assignmentScope[managementv1.ProjectAssignment] {
	api := client.PermissionsOpenfgaAPI
	return assignmentScope[managementv1.ProjectAssignment]{
		get: func(ctx context.Context) ([]managementv1.ProjectAssignment, error) {
			resp, _, err := api.GetProjectAssignmentsById(ctx, projectID).Execute()
			if err != nil {
				return nil, apiError(err)
			}
			return resp.Assignments, nil
		},
		update: func(ctx context.Context, writes, deletes []managementv1.ProjectAssignment) error {
			_, err := api.UpdateProjectAssignmentsById(ctx, projectID).UpdateProjectAssignmentsRequest(managementv1.UpdateProjectAssignmentsRequest{
				Writes:  writes,
				Deletes: deletes,
			}).Execute()
			return apiError(err)
		},
	}
}

func warehouseAssignmentScope(client *lakekeeper.Client, warehouseID string) assignmentScope[managementv1.WarehouseAssignment] {
	api := client.PermissionsOpenfgaAPI
	return assignmentScope[managementv1.WarehouseAssignment]{
		get: func(ctx context.Context) ([]managementv1.WarehouseAssignment, error) {
			resp, _, err := api.GetWarehouseAssignmentsById(ctx, warehouseID).Execute()
			if err != nil {
				return nil, apiError(err)
			}
			return resp.Assignments, nil
		},
		update: func(ctx context.Context, writes, deletes []managementv1.WarehouseAssignment) error {
			_, err := api.UpdateWarehouseAssignmentsById(ctx, warehouseID).UpdateWarehouseAssignmentsRequest(managementv1.UpdateWarehouseAssignmentsRequest{
				Writes:  writes,
				Deletes: deletes,
			}).Execute()
			return apiError(err)
		},
	}
}

func roleAssignmentScope(client *lakekeeper.Client, roleID string) assignmentScope[managementv1.RoleAssignment] {
	api := client.PermissionsOpenfgaAPI
	return assignmentScope[managementv1.RoleAssignment]{
		get: func(ctx context.Context) ([]managementv1.RoleAssignment, error) {
			resp, _, err := api.GetRoleAssignmentsById(ctx, roleID).Execute()
			if err != nil {
				return nil, apiError(err)
			}
			return resp.Assignments, nil
		},
		update: func(ctx context.Context, writes, deletes []managementv1.RoleAssignment) error {
			_, err := api.UpdateRoleAssignmentsById(ctx, roleID).UpdateRoleAssignmentsRequest(managementv1.UpdateRoleAssignmentsRequest{
				Writes:  writes,
				Deletes: deletes,
			}).Execute()
			return apiError(err)
		},
	}
}

func namespaceAssignmentScope(client *lakekeeper.Client, namespaceID string) assignmentScope[managementv1.NamespaceAssignment] {
	api := client.PermissionsOpenfgaAPI
	return assignmentScope[managementv1.NamespaceAssignment]{
		get: func(ctx context.Context) ([]managementv1.NamespaceAssignment, error) {
			resp, _, err := api.GetNamespaceAssignmentsById(ctx, namespaceID).Execute()
			if err != nil {
				return nil, apiError(err)
			}
			return resp.Assignments, nil
		},
		update: func(ctx context.Context, writes, deletes []managementv1.NamespaceAssignment) error {
			_, err := api.UpdateNamespaceAssignmentsById(ctx, namespaceID).UpdateNamespaceAssignmentsRequest(managementv1.UpdateNamespaceAssignmentsRequest{
				Writes:  writes,
				Deletes: deletes,
			}).Execute()
			return apiError(err)
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

// apply grants the added relations to one principal and revokes the removed
// ones, in a single request. It sends nothing when both lists are empty.
func (s assignmentScope[T]) apply(ctx context.Context, kind permissions.PrincipalKind, principalID string, added, removed []string) error {
	if len(added) == 0 && len(removed) == 0 {
		return nil
	}

	writes, err := buildAssignments[T](added, kind, principalID)
	if err != nil {
		return err
	}

	deletes, err := buildAssignments[T](removed, kind, principalID)
	if err != nil {
		return err
	}

	return s.update(ctx, writes, deletes)
}

// buildAssignments returns nil for no relations, so that the empty side of an
// update is omitted from the request instead of being sent as an empty list.
func buildAssignments[T any](relations []string, kind permissions.PrincipalKind, principalID string) ([]T, error) {
	if len(relations) == 0 {
		return nil, nil
	}

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
