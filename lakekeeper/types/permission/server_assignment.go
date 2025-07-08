package types

import (
	"encoding/json"
	"fmt"
)

// ServerAssignment represents an assignment a role or a user can
// have to the server
//
// Assignee can be a role or a user
// Assignement can be Operator or Admin
type ServerAssignment struct {
	Assignee   UserOrRole
	Assignment ServerAssignmentType
}

// to be sure ServerAssignment can be JSON encoded/decoded
var (
	_ json.Unmarshaler = (*ServerAssignment)(nil)
	_ json.Marshaler   = (*ServerAssignment)(nil)
)

type (
	UserOrRoleType string

	ServerAssignmentType string
)

const (
	UserType UserOrRoleType = "user"
	RoleType UserOrRoleType = "role"

	OperatorServerAssignment ServerAssignmentType = "operator"
	AdminServerAssignment    ServerAssignmentType = "admin"
)

type UserOrRole struct {
	Type  UserOrRoleType
	Value string
}

func (sa *ServerAssignment) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Type ServerAssignmentType `json:"type"`
		Role *string              `json:"role,omitempty"`
		User *string              `json:"user,omitempty"`
	}{}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	sa.Assignment = aux.Type

	if aux.Role == nil && aux.User == nil {
		return fmt.Errorf("error reading server assignment, role or user must be provided")
	}

	if aux.Role != nil && aux.User != nil {
		return fmt.Errorf("error reading server assignment, role and user can't be both provided")
	}

	if aux.Role != nil {
		sa.Assignee = UserOrRole{
			RoleType,
			*aux.Role,
		}
		return nil
	}

	if aux.User != nil {
		sa.Assignee = UserOrRole{
			UserType,
			*aux.User,
		}
		return nil
	}
	return fmt.Errorf("incorrect server assignment")
}

func (sa ServerAssignment) MarshalJSON() ([]byte, error) {
	aux := make(map[string]string)

	switch sa.Assignee.Type {
	case RoleType:
		aux["role"] = sa.Assignee.Value
	case UserType:
		aux["user"] = sa.Assignee.Value
	}

	aux["type"] = string(sa.Assignment)

	return json.Marshal(aux)
}
