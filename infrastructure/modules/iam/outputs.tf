
output "lambda_function_role_id" {
  value = "${aws_iam_role.lambda_function.arn}"
}

output "api_gateway_role_id" {
  value = "${aws_iam_role.gateway_invocation_role.arn}"
}