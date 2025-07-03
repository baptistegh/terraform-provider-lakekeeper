package lakekeeper

import (
	"encoding/json"
	"fmt"
)

var ValidDeleteProfileTypes = []string{"soft", "hard"}

type DeleteProfile interface {
	GetDeleteProfileType() string
}

type SoftDeleteProfile struct {
	Type           string `json:"type"`
	ExpiredSeconds int32  `json:"expired-seconds"`
}

func (SoftDeleteProfile) GetDeleteProfileType() string {
	return "soft"
}

type HardDeleteProfile struct {
	Type string `json:"type"`
}

func (HardDeleteProfile) GetDeleteProfileType() string {
	return "hard"
}

type DeleteProfileWrapper struct {
	DeleteProfile DeleteProfile
}

func (w *DeleteProfileWrapper) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	t, _ := raw["type"].(string)
	switch t {
	case "soft":
		var s SoftDeleteProfile
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.DeleteProfile = s
	case "hard":
		var s HardDeleteProfile
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.DeleteProfile = s
	default:
		return fmt.Errorf("unknown delete-profile type: %s", t)
	}
	return nil
}

func (w DeleteProfileWrapper) MarshalJSON() ([]byte, error) {
	if w.DeleteProfile == nil {
		return nil, nil
	}
	return json.Marshal(w.DeleteProfile)
}
