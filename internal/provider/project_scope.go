package provider

import (
	"errors"

	lakekeeper "github.com/lakekeeper/go-lakekeeper/pkg/client"
)

// errNoProject is returned by inProject for an empty project ID.
var errNoProject = errors.New("project ID must not be empty")

// inProject scopes a project-scoped façade call to projectID.
//
// A façade call without a project runs in the server's default project, so an
// empty ID would silently read, change or delete another project. inProject
// rejects it instead. Use a call without any option to target the default
// project on purpose.
func inProject(projectID string) (lakekeeper.CallOption, error) {
	if projectID == "" {
		return nil, errNoProject
	}

	return lakekeeper.WithProject(projectID), nil
}
