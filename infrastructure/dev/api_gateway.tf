
resource "aws_api_gateway_rest_api" "api" {
  name = "${var.project}-${var.apex_environment}-api"
  description = "${var.apex_environment} serverless API for ${var.project}"
}

resource "aws_api_gateway_resource" "v1" {
  rest_api_id = "${aws_api_gateway_rest_api.api.id}"
  parent_id = "${aws_api_gateway_rest_api.api.root_resource_id}"
  path_part = "v1"
}

resource "aws_api_gateway_deployment" "dev" {
  depends_on = [
    "aws_api_gateway_integration.PostSlackStandupToLambda",
  ]
  rest_api_id = "${aws_api_gateway_rest_api.api.id}"
  stage_name = "${var.apex_environment}"
  description = "${var.apex_environment} deployment"
}