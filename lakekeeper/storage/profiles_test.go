package storage

import (
	"encoding/json"
	"testing"
)

func TestStorageProfileWrapper_S3_MarshalJSON(t *testing.T) {
	given := StorageProfileWrapper{
		StorageProfile: StorageProfileS3{
			Type:       "s3",
			Bucket:     "test-bucket",
			Region:     "us-west-2",
			STSEnabled: false,
		},
	}

	r, err := json.Marshal(given)
	if err != nil {
		t.Fatalf("failed to marshal StorageProfileWrapper: %v", err)
	}

	expected := `{"type":"s3","bucket":"test-bucket","region":"us-west-2","sts-enabled":false}`

	if string(r) != expected {
		t.Errorf("expected %s, got %s", expected, string(r))
	}
}

func TestStorageProfileWrapper_GCS_MarshalJSON(t *testing.T) {
	given := StorageProfileWrapper{
		StorageProfile: GCSStorageSettings{
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
	var given StorageProfileWrapper
	if err := json.Unmarshal([]byte(expected), &given); err != nil {
		t.Fatalf("failed to unmarshal StorageProfileWrapper, %v", err)
	}

	switch given.StorageProfile.(type) {
	case ADLSStorageSettings:
		// expected type
	default:
		t.Errorf("expected ADLSStorageSettings, got %T", given.StorageProfile)
	}
}

func TestStorageProfileWrapper_S3_UnmarshalJSON(t *testing.T) {
	expected := `{"type":"s3","bucket":"test-bucket","region":"us-west-2","sts-enabled":false}`
	var given StorageProfileWrapper
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
	var given StorageProfileWrapper
	if err := json.Unmarshal([]byte(expected), &given); err != nil {
		t.Fatalf("failed to unmarshal StorageProfileWrapper, %v", err)
	}

	switch given.StorageProfile.(type) {
	case GCSStorageSettings:
		// expected type
	default:
		t.Errorf("expected GCSStorageSettings, got %T", given.StorageProfile)
	}
}
