package ecxfabric

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"ecxfabric": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

type preCheckFunc = func(*testing.T)

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ECXFABRIC_CLIENT_ID"); v == "" {
		t.Fatal("ECXFABRIC_CLIENTID must be set for acceptance tests")
	}

	if v := os.Getenv("ECXFABRIC_CLIENT_SECRET"); v == "" {
		t.Fatal("ECXFABRIC_CLIENT_SECRET must be set for acceptance tests")
	}

	if v := os.Getenv("ECXFABRIC_USERNAME"); v == "" {
		t.Fatal("ECXFABRIC_USERNAME must be set for acceptance tests")
	}

	if v := os.Getenv("ECXFABRIC_PASSWORD"); v == "" {
		t.Fatal("ECXFABRIC_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("AWS_ACCESS_KEY_ID"); v == "" {
		t.Fatal("AWS_ACCESS_KEY_ID must be set for acceptance tests")
	}

	if v := os.Getenv("AWS_SECRET_ACCESS_KEY"); v == "" {
		t.Fatal("AWS_SECRET_ACCESS_KEY must be set for acceptance tests")
	}
}
