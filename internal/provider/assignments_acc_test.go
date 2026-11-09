//go:build acceptance

package provider

import (
	"context"

	"github.com/lakekeeper/go-lakekeeper/pkg/permissions"
)

// testAssignment grants relation to one principal.
type testAssignment struct {
	kind        permissions.PrincipalKind
	principalID string
	relation    string
}

// grantAssignments writes the given assignments on the object of scope, in one request.
func grantAssignments[T any](ctx context.Context, scope assignmentScope[T], assignments ...testAssignment) error {
	writes := make([]T, 0, len(assignments))
	for _, a := range assignments {
		w, err := permissions.BuildAssignment[T](a.relation, a.kind, a.principalID)
		if err != nil {
			return err
		}
		writes = append(writes, w)
	}

	return scope.grant(ctx, writes)
}
