package types

import (
	"encoding/json"
	"testing"
)

func TestServerPermissions_MarshalJSON(t *testing.T) {
	expected := []string{
		`{"role":"a6e5a780-258e-4bee-9bd8-f8ae3f675415","type":"admin"}`,
		`{"role":"9cc096bf-db1f-43f3-bea6-f0819df32db0","type":"operator"}`,
		`{"type":"admin","user":"f5c2329c-8679-44d0-8ea3-167ee14fa94e"}`,
		`{"type":"operator","user":"a0d21f3d-2cbb-4066-8b77-5ec5a21680be"}`,
	}

	given := []ServerAssignment{
		{
			Assignment: AdminServerAssignment,
			Assignee: UserOrRole{
				Type:  RoleType,
				Value: "a6e5a780-258e-4bee-9bd8-f8ae3f675415",
			},
		},
		{
			Assignment: OperatorServerAssignment,
			Assignee: UserOrRole{
				Type:  RoleType,
				Value: "9cc096bf-db1f-43f3-bea6-f0819df32db0",
			},
		},
		{
			Assignment: AdminServerAssignment,
			Assignee: UserOrRole{
				Type:  UserType,
				Value: "f5c2329c-8679-44d0-8ea3-167ee14fa94e",
			},
		},
		{
			Assignment: OperatorServerAssignment,
			Assignee: UserOrRole{
				Type:  UserType,
				Value: "a0d21f3d-2cbb-4066-8b77-5ec5a21680be",
			},
		},
	}

	for k, v := range expected {
		b, err := json.Marshal(given[k])
		if err != nil {
			t.Fatalf("%v", err)
		}
		if string(b) != v {
			t.Fatalf("exepected %s got %s", v, string(b))
		}
	}
}

func TestServerPermissions_UnmarshalJSON(t *testing.T) {
	expected := []ServerAssignment{
		{
			Assignment: AdminServerAssignment,
			Assignee: UserOrRole{
				Type:  RoleType,
				Value: "a6e5a780-258e-4bee-9bd8-f8ae3f675415",
			},
		},
		{
			Assignment: OperatorServerAssignment,
			Assignee: UserOrRole{
				Type:  RoleType,
				Value: "9cc096bf-db1f-43f3-bea6-f0819df32db0",
			},
		},
		{
			Assignment: AdminServerAssignment,
			Assignee: UserOrRole{
				Type:  UserType,
				Value: "f5c2329c-8679-44d0-8ea3-167ee14fa94e",
			},
		},
		{
			Assignment: OperatorServerAssignment,
			Assignee: UserOrRole{
				Type:  UserType,
				Value: "a0d21f3d-2cbb-4066-8b77-5ec5a21680be",
			},
		},
	}

	given := []string{
		`{"role":"a6e5a780-258e-4bee-9bd8-f8ae3f675415","type":"admin"}`,
		`{"role":"9cc096bf-db1f-43f3-bea6-f0819df32db0","type":"operator"}`,
		`{"type":"admin","user":"f5c2329c-8679-44d0-8ea3-167ee14fa94e"}`,
		`{"type":"operator","user":"a0d21f3d-2cbb-4066-8b77-5ec5a21680be"}`,
	}

	for k, v := range expected {
		var aux ServerAssignment
		err := json.Unmarshal([]byte(given[k]), &aux)
		if err != nil {
			t.Fatalf("%v", err)
		}

		if v.Assignment != aux.Assignment {
			t.Fatalf("expected %s got %s", v.Assignment, aux.Assignment)
		}

		if v.Assignee.Type != aux.Assignee.Type {
			t.Fatalf("expected %s got %s", v.Assignee.Type, aux.Assignee.Type)
		}

		if v.Assignee.Value != aux.Assignee.Value {
			t.Fatalf("expected %s got %s", v.Assignee.Type, aux.Assignee.Value)
		}
	}
}
