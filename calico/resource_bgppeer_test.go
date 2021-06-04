package calico

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/projectcalico/libcalico-go/lib/options"
	"strings"
	"testing"
)

func TestAccBgpPeer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBgpPeerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateBgpPeerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBgpPeerExists("calico_bgppeer.test"),
					resource.TestCheckResourceAttr("calico_bgppeer.test", "metadata.0.name", "testbgppeer"),
					resource.TestCheckResourceAttr("calico_bgppeer.test", "spec.0.as_number", "62512"),
					resource.TestCheckResourceAttr("calico_bgppeer.test", "spec.0.node_selector", "k == 'v'"),
					resource.TestCheckResourceAttr("calico_bgppeer.test", "spec.0.peer_ip", "192.168.0.5"),
				),
			},
			{
				Config: testUpdateBgpPeerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBgpPeerExists("calico_bgppeer.test"),
					resource.TestCheckResourceAttr("calico_bgppeer.test", "metadata.0.name", "testbgppeer2"),
					resource.TestCheckResourceAttr("calico_bgppeer.test", "spec.0.node", "global"),
					resource.TestCheckResourceAttr("calico_bgppeer.test", "spec.0.peer_selector", "k == 'v'"),
				),
			},
		},
	})
}

func testAccCheckBgpPeerDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(config).Client

	apis := getResourcesByType("calico_bgppeer", state)

	if len(apis) != 1 {
		return fmt.Errorf("expecting only 1 BgpPeer resource found %v", len(apis))
	}

	_, err := client.BGPPeers().Get(ctx, apis[0].Primary.ID, options.GetOptions{})

	switch {
	case err == nil:
		return errors.New("Expected error, got none")
	case err != nil && !strings.Contains(err.Error(), "resource does not exist"):
		return fmt.Errorf("Expected 404, got %s", err)
	}

	return nil
}

func testAccCheckBgpPeerExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		BgpPeer, err := testAccProvider.Meta().(config).Client.BGPPeers().Get(ctx, rs.Primary.ID, options.GetOptions{})

		if err != nil {
			return err
		}

		if BgpPeer == nil {
			return fmt.Errorf("BgpPeer with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

const testCreateBgpPeerConfig = `
resource "calico_bgppeer" "test" {
  metadata{
    name = "testbgppeer"
  }
  spec{
	as_number = 62512
	node_selector = "k == 'v'"
    peer_ip = "192.168.0.5"
  }
}`

const testUpdateBgpPeerConfig = `
resource "calico_bgppeer" "test" {
  metadata{
    name = "testbgppeer2"
  }
  spec{
    node = "global"
	peer_selector = "k == 'v'"
  }
}`
