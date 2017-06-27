
output "lambda_function_role_id" {
  value = "${module.iam.lambda_function_role_id}"
}

output "api_gateway_role_id" {
  value = "${module.iam.api_gateway_role_id}"
}

output "lambda_security_group_id" {
  value = "${aws_security_group.lambda_security_group.id}"
}

output "private_subnets" {
  value = ["${module.vpc.private_subnets}"]
}

output "public_subnets" {
  value = ["${module.vpc.public_subnets}"]
}

output "vpc_id" {
  value = "${module.vpc.vpc_id}"
}

output "dynamodb_status_events_tablename" {
  value = "${aws_dynamodb_table.status_events_table.id}"
}

output "dynamodb_status_aggregates_tablename" {
  value = "${aws_dynamodb_table.status_aggregates_table.id}"
}
