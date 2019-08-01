package ecxfabric

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/nomical/terraform-provider-ecxfabric/apiclient"
)

func TestAccL2Connection_Basic(t *testing.T) {
	resourceID := "id" + acctest.RandString(10) // Added "id" at the beggining as RandString generated integer at start of string (sometimes!). TF Error: Invalid resource name; A name must start with a letter and may contain only letters, digits, underscores, and dashes.
	authorizationKey := "986744318870"
	notifications1 := "support@domain.com"
	primaryName := "TF-L2C-Test"
	purchaseOrderNumber := "TF-L2C-Test-PO"
	primaryPortUUID := "7b5650d1-810a-10a0-66e0-30ac094f8701"
	primaryVlanSTag := "20"
	profileUUID := "69ee618d-be52-468d-bc99-00566f2dd2b9"
	sellerMetroCode := "LD"
	sellerRegion := "eu-west-2"
	speed := "50"
	speedUnit := "MB"
	name := "ecxfabric_l2_connection." + resourceID

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckL2ConnectionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testL2ConnectionConfig(resourceID,
					authorizationKey,
					notifications1,
					primaryName,
					purchaseOrderNumber,
					primaryPortUUID,
					primaryVlanSTag,
					profileUUID,
					sellerMetroCode,
					sellerRegion,
					speed,
					speedUnit),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "authorization_key", authorizationKey),
					resource.TestCheckResourceAttr(name, "notifications.0", notifications1),
					resource.TestCheckResourceAttr(name, "primary_name", primaryName),
					resource.TestCheckResourceAttr(name, "purchase_order_number", purchaseOrderNumber),
					resource.TestCheckResourceAttr(name, "primary_port_uuid", primaryPortUUID),
					resource.TestCheckResourceAttr(name, "primary_vlan_s_tag", primaryVlanSTag),
					resource.TestCheckResourceAttr(name, "profile_uuid", profileUUID),
					resource.TestCheckResourceAttr(name, "seller_metro_code", sellerMetroCode),
					resource.TestCheckResourceAttr(name, "seller_region", sellerRegion),
					resource.TestCheckResourceAttr(name, "speed", speed),
					resource.TestCheckResourceAttr(name, "speed_unit", speedUnit),
				),
			},
		},
	})
}

func testAccCheckL2ConnectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecxfabric_l2_connection" {
			continue
		}

		conn, err := client.ReadL2Connection(rs.Primary.ID)
		if err != nil {
			return err
		}

		if conn.Status != apiclient.L2ConnectionProviderStatusDeprovisioning {
			return fmt.Errorf("L2 connection in unexpected state! UUID: %v, Expected: %v, Actual: %v", rs.Primary.ID, apiclient.L2ConnectionProviderStatusDeprovisioning, conn.Status)
		}
	}

	return nil
}

func testL2ConnectionConfig(resourceID, authorizationKey, notifications1, primaryName, purchaseOrderNumber, primaryPortUUID, primaryVlanSTag, profileUUID, sellerMetroCode, sellerRegion, speed, speedUnit string) string {
	return fmt.Sprintf(`
	resource "ecxfabric_l2_connection" "%[1]s" {
		authorization_key    = "%[2]s"
		notifications       = ["%[3]s"]
		primary_name         = "%[4]s"
		purchase_order_number = "%[5]s"
		primary_port_uuid     = "%[6]s"
		primary_vlan_s_tag     = %[7]s
		profile_uuid       = "%[8]s"
		seller_metro_code     = "%[9]s"
		seller_region        = "%[10]s"
		speed               = %[11]s
		speed_unit           = "%[12]s"
	}
	`, resourceID, authorizationKey, notifications1, primaryName, purchaseOrderNumber, primaryPortUUID, primaryVlanSTag, profileUUID, sellerMetroCode, sellerRegion, speed, speedUnit)
}
