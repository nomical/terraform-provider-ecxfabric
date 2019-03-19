output "ecxfabric_l2_connection" {
  value = <<ALL
        aws_connection_id = ${ecxfabric_l2_connection.test.aws_connection_id}
    ALL
}

output "ecxfabric_l2_connection_aws_accepter" {
  value = <<ALL
        aws_device = ${ecxfabric_l2_connection_aws_accepter.accepter.aws_device}
        aws_device_v2 = ${ecxfabric_l2_connection_aws_accepter.accepter.aws_device_v2}
        bandwidth = ${ecxfabric_l2_connection_aws_accepter.accepter.bandwidth}
        connection_id = ${ecxfabric_l2_connection_aws_accepter.accepter.connection_id}
        connection_name = ${ecxfabric_l2_connection_aws_accepter.accepter.connection_name}
        connection_state = ${ecxfabric_l2_connection_aws_accepter.accepter.connection_state}
        jumbo_frame_capable = ${ecxfabric_l2_connection_aws_accepter.accepter.jumbo_frame_capable}
        location = ${ecxfabric_l2_connection_aws_accepter.accepter.location}
        owner_account = ${ecxfabric_l2_connection_aws_accepter.accepter.owner_account}
        partner_name = ${ecxfabric_l2_connection_aws_accepter.accepter.partner_name}
        region = ${ecxfabric_l2_connection_aws_accepter.accepter.region}
        vlan = ${ecxfabric_l2_connection_aws_accepter.accepter.vlan}
    ALL
}
