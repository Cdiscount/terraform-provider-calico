package calico

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]terraform.ResourceProvider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"calico": testAccProvider,
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

func TestAccPreCheck(t *testing.T) {
	variables := []string{
		"CALICO_BACKEND_TYPE",
		"CALICO_ETCD_ENDPOINTS",
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Fatalf("`%s` must be set for acceptance tests!", variable)
		}
	}
}

func TestMain(m *testing.M) {
	code := m.Run()

	os.Exit(code)
}
