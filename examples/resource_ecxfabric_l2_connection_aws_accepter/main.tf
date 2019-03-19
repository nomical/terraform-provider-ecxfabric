resource "ecxfabric_l2_connection" "test" {
  authorization_key     = "${var.aws_account_id}"
  notifications         = "${var.notifications}"
  primary_name          = "${var.primary_name}"
  purchase_order_number = "${var.purchase_order_number}"
  primary_port_uuid     = "${var.primary_port_uuid}"
  primary_vlan_s_tag    = "${var.primary_vlan_s_tag}"
  profile_uuid          = "${var.profile_uuid}"
  seller_metro_code     = "${var.seller_metro_code}"
  seller_region         = "${var.seller_region}"
  speed                 = "${var.speed}"
  speed_unit            = "${var.speed_unit}"
}

resource "ecxfabric_l2_connection_aws_accepter" "accepter" {
  access_key    = "${var.aws_access_key}"
  secret_key    = "${var.aws_secret_key}"
  region        = "${var.aws_region}"
  connection_id = "${ecxfabric_l2_connection.test.aws_connection_id}"
}
