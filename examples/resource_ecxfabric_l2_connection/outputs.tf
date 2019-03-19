output "ecxfabric_l2_connection" {
  value = <<ALL
        aws_connection_id = ${ecxfabric_l2_connection.test.aws_connection_id}
    ALL
}
