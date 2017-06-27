# api resource for slack slash commands
resource "aws_api_gateway_resource" "Slack" {
  rest_api_id = "${aws_api_gateway_rest_api.api.id}"
  parent_id = "${aws_api_gateway_resource.v1.id}"
  path_part = "slack"
}

# api resource for slack slash command `/standup`
resource "aws_api_gateway_resource" "SlackStandup" {
  rest_api_id = "${aws_api_gateway_rest_api.api.id}"
  parent_id = "${aws_api_gateway_resource.Slack.id}"
  path_part = "standup"
}

# CORS settings for `/slack/standup`
//module "standup_cors" {
//  source      = "../modules/cors"
//  rest_api_id = "${aws_api_gateway_rest_api.api.id}"
//  resource_id = "${aws_api_gateway_resource.SlackStandup.id}"
//}

resource "aws_api_gateway_method" "PostSlackStandup" {
  rest_api_id = "${aws_api_gateway_rest_api.api.id}"
  resource_id = "${aws_api_gateway_resource.SlackStandup.id}"
  http_method = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "PostSlackStandupToLambda" {
  rest_api_id = "${aws_api_gateway_rest_api.api.id}"
  resource_id = "${aws_api_gateway_resource.SlackStandup.id}"
  http_method = "${aws_api_gateway_method.PostSlackStandup.http_method}"
  type = "AWS_PROXY"
  integration_http_method = "POST"
  credentials = "${module.iam.api_gateway_role_id}"
  # http://docs.aws.amazon.com/apigateway/api-reference/resource/integration/#uri
  uri = "arn:aws:apigateway:${var.aws_region}:lambda:path/2015-03-31/functions/${var.apex_function_slack_standup}/invocations"
}
