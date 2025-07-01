package lakekeeper

import (
	"encoding/json"
	"testing"
)

func TestStorageProfileWrapper_S3_MarshalJSON(t *testing.T) {
	given := storageProfileWrapper{
		StorageProfile: StorageProfileS3{
			Type:       "s3",
			Bucket:     "test-bucket",
			Region:     "us-west-2",
			STSEnabled: false,
		},
	}

	r, err := json.Marshal(given)
	if err != nil {
		t.Fatalf("failed to marshal storageProfileWrapper: %v", err)
	}

	expected := `{"type":"s3","bucket":"test-bucket","region":"us-west-2","sts-enabled":false}`

	if string(r) != expected {
		t.Errorf("expected %s, got %s", expected, string(r))
	}
}

func TestStorageProfileWrapper_GCS_MarshalJSON(t *testing.T) {
	given := storageProfileWrapper{
		StorageProfile: StorageProfileGCS{
			Type:   "gcs",
			Bucket: "test-bucket",
		},
	}

	r, err := json.Marshal(given)
	if err != nil {
		t.Fatalf("failed to marshal Warehouse: %v", err)
	}

	expected := `{"type":"gcs","bucket":"test-bucket"}`

	if string(r) != expected {
		t.Errorf("expected %s, got %s", expected, string(r))
	}
}

func TestStorageProfileWrapper_ADLS_UnmarshalJSON(t *testing.T) {
	expected := `{"type":"adls","account-name":"test-account","allow-alternative-protocols":false,"authority-host":"","filesystem":"test-filesystem","host":"","key-prefix":"","sas-token-validity-seconds":0}`
	var given storageProfileWrapper
	if err := json.Unmarshal([]byte(expected), &given); err != nil {
		t.Fatalf("failed to unmarshal storageProfileWrapper, %v", err)
	}

	switch given.StorageProfile.(type) {
	case StorageProfileADLS:
		// expected type
	default:
		t.Errorf("expected StorageProfileADLS, got %T", given.StorageProfile)
	}
}

func TestStorageProfileWrapper_S3_UnmarshalJSON(t *testing.T) {
	expected := `{"type":"s3","bucket":"test-bucket","region":"us-west-2","sts-enabled":false}`
	var given storageProfileWrapper
	if err := json.Unmarshal([]byte(expected), &given); err != nil {
		t.Fatalf("failed to unmarshal Warehouse: %v", err)
	}

	switch given.StorageProfile.(type) {
	case StorageProfileS3:
		// expected type
	default:
		t.Errorf("expected StorageProfileS3, got %T", given.StorageProfile)
	}
}

func TestStorageProfileWrapper_GCS_UnmarshalJSON(t *testing.T) {
	expected := `{"type":"gcs","bucket":"test-bucket"}`
	var given storageProfileWrapper
	if err := json.Unmarshal([]byte(expected), &given); err != nil {
		t.Fatalf("failed to unmarshal storageProfileWrapper, %v", err)
	}

	switch given.StorageProfile.(type) {
	case StorageProfileGCS:
		// expected type
	default:
		t.Errorf("expected StorageProfileGCS, got %T", given.StorageProfile)
	}
}
