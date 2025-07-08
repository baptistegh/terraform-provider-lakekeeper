package profile

import (
	"encoding/json"
	"fmt"
)

type DeleteProfileSettings interface {
	GetDeteProfileType() DeleteProfileType
	AsProfile() *DeleteProfile

	json.Marshaler
}

type DeleteProfile struct {
	DeleteProfileSettings DeleteProfileSettings
}

type DeleteProfileType string

const (
	HardDeleteProfileType DeleteProfileType = "hard"
	SoftDeleteProfileType DeleteProfileType = "soft"
)

var (
	_ DeleteProfileSettings = (*TabularDeleteProfileHard)(nil)
	_ DeleteProfileSettings = (*TabularDeleteProfileSoft)(nil)
)

type TabularDeleteProfileHard struct{}

func NewTabularDeleteProfileHard() *TabularDeleteProfileHard {
	return &TabularDeleteProfileHard{}
}

func (*TabularDeleteProfileHard) GetDeteProfileType() DeleteProfileType {
	return HardDeleteProfileType
}

func (d *TabularDeleteProfileHard) AsProfile() *DeleteProfile {
	return &DeleteProfile{DeleteProfileSettings: d}
}

func (d TabularDeleteProfileHard) MarshalJSON() ([]byte, error) {
	aux := struct {
		Type string `json:"type"`
	}{
		Type: string(d.GetDeteProfileType()),
	}
	return json.Marshal(aux)
}

type TabularDeleteProfileSoft struct {
	ExpirationSeconds int32 `json:"expiration-seconds"`
}

func NewTabularDeleteProfileSoft(expirationSeconds int32) *TabularDeleteProfileSoft {
	return &TabularDeleteProfileSoft{
		ExpirationSeconds: expirationSeconds,
	}
}

func (*TabularDeleteProfileSoft) GetDeteProfileType() DeleteProfileType {
	return SoftDeleteProfileType
}

func (d *TabularDeleteProfileSoft) AsProfile() *DeleteProfile {
	return &DeleteProfile{DeleteProfileSettings: d}
}

func (d TabularDeleteProfileSoft) MarshalJSON() ([]byte, error) {
	type Alias TabularDeleteProfileSoft
	aux := struct {
		Type string `json:"type"`
		Alias
	}{
		Type:  string(d.GetDeteProfileType()),
		Alias: Alias(d),
	}
	return json.Marshal(aux)
}

func (sc *DeleteProfile) UnmarshalJSON(data []byte) error {
	var peek struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	switch peek.Type {
	case "hard":
		var cfg TabularDeleteProfileHard
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.DeleteProfileSettings = &cfg
	case "soft":
		var cfg TabularDeleteProfileSoft
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.DeleteProfileSettings = &cfg
	default:
		return fmt.Errorf("unsupported delete profile type: %s", peek.Type)
	}
	return nil
}

func (sc DeleteProfile) MarshalJSON() ([]byte, error) {
	return json.Marshal(sc.DeleteProfileSettings)
}
