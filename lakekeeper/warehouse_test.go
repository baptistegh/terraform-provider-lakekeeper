package lakekeeper

import (
	"encoding/json"
	"testing"

	"github.com/baptistegh/terraform-provider-lakekeeper/lakekeeper/storage"
)

func TestCreateWarehouseOptions_Marshal(t *testing.T) {
	opts := CreateWarehouseOptions{
		Name:              "name",
		ProjectID:         "prj_id",
		StorageProfile:    storage.StorageProfileWrapper{StorageProfile: storage.NewStorageProfileS3("test", "us-west-1", true)},
		StorageCredential: storage.StorageCredentialWrapper{StorageCredential: storage.NewS3CredentialAccessKey("keyid", "secretkey", "")},
	}

	b, err := json.Marshal(opts)
	if err != nil {
		t.Fatalf("%v", err)
	}

	expected := `{"warehouse-name":"name","project-id":"prj_id","storage-profile":{"type":"s3","bucket":"test","region":"us-west-1","sts-enabled":true},"storage-credential":{"type":"s3","credential-type":"access-key","aws-access-key-id":"keyid","aws-secret-access-key":"secretkey"}}`

	if string(b) != expected {
		t.Fatalf("expected %s got %s", expected, string(b))
	}
}
