package ecxfabric

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/directconnect"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccL2ConnectionAwsAccepter_Basic(t *testing.T) {
	connID := "id" + acctest.RandString(10)     // Added "id" at the beggining as RandString generated integer at start of string (sometimes!). TF Error: Invalid resource name; A name must start with a letter and may contain only letters, digits, underscores, and dashes.
	accepterID := "id" + acctest.RandString(10) // Added "id" at the beggining as RandString generated integer at start of string (sometimes!). TF Error: Invalid resource name; A name must start with a letter and may contain only letters, digits, underscores, and dashes.
	accepterName := "ecxfabric_l2_connection_aws_accepter." + accepterID

	authorizationKey := "986744318870"
	notifications1 := "support@domain.com"
	primaryName := "TF-L2CAwsAccepter-Test"
	purchaseOrderNumber := "TF-L2CAwsAccepter-Test-PO"
	primaryPortUUID := "7b5650d1-810a-10a0-66e0-30ac094f8701"
	primaryVlanSTag := "20"
	profileUUID := "69ee618d-be52-468d-bc99-00566f2dd2b9"
	sellerMetroCode := "LD"
	sellerRegion := "eu-west-2"
	speed := "50"
	speedUnit := "MB"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckL2ConnectionAwsAccepterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testL2ConnectionAwsAccepterConfig(connID, authorizationKey, notifications1, primaryName,
					purchaseOrderNumber, primaryPortUUID, primaryVlanSTag, profileUUID, sellerMetroCode, sellerRegion,
					speed, speedUnit, accepterID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(accepterName, "connection_state", directconnect.ConnectionStateAvailable),
				),
			},
		},
	})
}

func testAccCheckL2ConnectionAwsAccepterDestroy(s *terraform.State) error {
	return nil
}

func testL2ConnectionAwsAccepterConfig(connID, authorizationKey, notifications1, primaryName,
	purchaseOrderNumber, primaryPortUUID, primaryVlanSTag, profileUUID, sellerMetroCode, sellerRegion,
	speed, speedUnit, accepterID string) string {
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

	resource "ecxfabric_l2_connection_aws_accepter" "%[13]s" {
	   connection_id = "${ecxfabric_l2_connection.%[1]s.aws_connection_id}"
	}
	`, connID, authorizationKey, notifications1, primaryName,
		purchaseOrderNumber, primaryPortUUID, primaryVlanSTag, profileUUID, sellerMetroCode, sellerRegion,
		speed, speedUnit, accepterID)
}
