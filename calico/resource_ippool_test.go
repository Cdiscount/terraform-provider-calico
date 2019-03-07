package calico

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/projectcalico/libcalico-go/lib/options"
	"strings"
	"testing"
)

func TestAccIpPool(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIpPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateIpPoolConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIpPoolExists("calico_ippool.test"),
					resource.TestCheckResourceAttr("calico_ippool.test", "metadata.0.name", "testippool"),
					resource.TestCheckResourceAttr("calico_ippool.test", "spec.0.block_size", "26"),
					resource.TestCheckResourceAttr("calico_ippool.test", "spec.0.cidr", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("calico_ippool.test", "spec.0.disabled", "true"),
					resource.TestCheckResourceAttr("calico_ippool.test", "spec.0.ipip_mode", "Never"),
					resource.TestCheckResourceAttr("calico_ippool.test", "spec.0.nat_outgoing", "true"),
					resource.TestCheckResourceAttr("calico_ippool.test", "spec.0.node_selector", "k == 'v'"),
				),
			},
			{
				Config: testUpdateIpPoolConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIpPoolExists("calico_ippool.test"),
					resource.TestCheckResourceAttr("calico_ippool.test", "metadata.0.name", "testippool2"),
					resource.TestCheckResourceAttr("calico_ippool.test", "spec.0.cidr", "192.168.0.0/24"),
					resource.TestCheckResourceAttr("calico_ippool.test", "spec.0.disabled", "false"),
					resource.TestCheckResourceAttr("calico_ippool.test", "spec.0.ipip_mode", "Always"),
					resource.TestCheckResourceAttr("calico_ippool.test", "spec.0.nat_outgoing", "false"),
				),
			},
		},
	})
}

func testAccCheckIpPoolDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(config).Client

	apis := getResourcesByType("calico_ippool", state)

	if len(apis) != 1 {
		return fmt.Errorf("expecting only 1 ippool resource found %v", len(apis))
	}

	_, err := client.IPPools().Get(ctx, apis[0].Primary.ID, options.GetOptions{})

	switch {
	case err == nil:
		return errors.New("Expected error, got none")
	case err != nil && !strings.Contains(err.Error(), "resource does not exist"):
		return fmt.Errorf("Expected 404, got %s", err)
	}

	return nil
}

func testAccCheckIpPoolExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ipPool, err := testAccProvider.Meta().(config).Client.IPPools().Get(ctx, rs.Primary.ID, options.GetOptions{})

		if err != nil {
			return err
		}

		if ipPool == nil {
			return fmt.Errorf("ippool with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

const testCreateIpPoolConfig = `
resource "calico_ippool" "test" {
  metadata{
    name = "testippool"
  }
  spec{
	block_size = "26"
    cidr = "192.168.1.0/24"
	disabled = true	
	ipip_mode = "Never"
    nat_outgoing = true
	node_selector = "k == 'v'"
  }
}`

const testUpdateIpPoolConfig = `
resource "calico_ippool" "test" {
  metadata{
    name = "testippool2"
  }
  spec{
    cidr = "192.168.0.0/24"
	disabled = false
	ipip_mode = "Always"
    nat_outgoing = false
  }
}`
