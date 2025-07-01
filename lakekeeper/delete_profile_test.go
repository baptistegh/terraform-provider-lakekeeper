package lakekeeper

import (
	"encoding/json"
	"testing"
)

func TestDeleteProfileWrapper_Soft_MarshalJSON(t *testing.T) {
	given := deleteProfileWrapper{
		DeleteProfile: SoftDeleteProfile{
			Type:           "soft",
			ExpiredSeconds: 3600,
		},
	}

	r, err := json.Marshal(given)
	if err != nil {
		t.Fatalf("failed to marshal deleteProfileWrapper: %v", err)
	}

	expected := `{"type":"soft","expired-seconds":3600}`

	if string(r) != expected {
		t.Errorf("expected %s, got %s", expected, string(r))
	}
}

func TestDeleteProfileWrapper_Hard_MarshalJSON(t *testing.T) {
	given := deleteProfileWrapper{
		DeleteProfile: HardDeleteProfile{
			Type: "hard",
		},
	}

	r, err := json.Marshal(given)
	if err != nil {
		t.Fatalf("failed to marshal Warehouse: %v", err)
	}

	expected := `{"type":"hard"}`

	if string(r) != expected {
		t.Errorf("expected %s, got %s", expected, string(r))
	}
}

func TestDeleteProfileWrapper_Soft_UnmarshalJSON(t *testing.T) {
	expected := `{"type":"soft","expired-seconds":3600}`
	var given deleteProfileWrapper
	if err := json.Unmarshal([]byte(expected), &given); err != nil {
		t.Fatalf("failed to unmarshal delete profile, %v", err)
	}

	switch given.DeleteProfile.(type) {
	case SoftDeleteProfile:
		// expected type
	default:
		t.Errorf("expected SoftDeleteProfile, got %T", given.DeleteProfile)
	}
}

func TestDeleteProfileWrapper_Hard_UnmarshalJSON(t *testing.T) {
	expected := `{"type":"hard"}`
	var given deleteProfileWrapper
	if err := json.Unmarshal([]byte(expected), &given); err != nil {
		t.Fatalf("failed to unmarshal deleteProfileWrapper, %v", err)
	}

	switch given.DeleteProfile.(type) {
	case HardDeleteProfile:
		// expected type
	default:
		t.Errorf("expected HardDeleteProfile, got %T", given.DeleteProfile)
	}
}
