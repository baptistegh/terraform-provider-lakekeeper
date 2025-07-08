package profile

import (
	"encoding/json"
	"testing"
)

func TestS3StorageProfile_Marshal(t *testing.T) {
	storageProfile, err := NewS3StorageSettings(
		"bucket1",
		"eu-west-1",
		WithEndpoint("http://minio:9000/"),
		WithPathStyleAccess(),
		WithS3KeyPrefix("warehouse"),
	)
	if err != nil {
		t.Fatalf("%v", err)
	}

	expected := `{"type":"s3","bucket":"bucket1","region":"eu-west-1","sts-enabled":false,"endpoint":"http://minio:9000/","key-prefix":"warehouse","path-style-access":true}`

	b, err := json.Marshal(storageProfile)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if string(b) != expected {
		t.Fatalf("expected %s got %s", expected, string(b))
	}

	// by Config
	b, err = json.Marshal(storageProfile.AsProfile())
	if err != nil {
		t.Fatalf("%v", err)
	}

	if string(b) != expected {
		t.Fatalf("expected %s got %s", expected, string(b))
	}
}
