#
# dynamodb table for status_events
#
resource "aws_dynamodb_table" "status_events_table" {
  name = "${var.project}-${var.apex_environment}-status-events"
  read_capacity = 10
  write_capacity = 10
  hash_key = "key"
  range_key = "partition"
  stream_enabled = "true"
  stream_view_type = "NEW_AND_OLD_IMAGES"

  attribute {
    name = "key"
    type = "S"
  }

  attribute {
    name = "partition"
    type = "N"
  }

  tags {
    Name = "${var.project}-${var.apex_environment}-dynamodb-table-status_events"
    Environment = "${var.apex_environment}"
    Project = "${var.project}"
  }
}

#
# dynamodb stream mapping for status_events aggregates
#
//resource "aws_lambda_event_source_mapping" "status_events_source_mapping" {
//  batch_size        = 1
//  event_source_arn  = "${aws_dynamodb_table.status_events_table.stream_arn}"
//  enabled           = true
//  function_name     = "${var.apex_function_aggregate_status_events}"
//  starting_position = "TRIM_HORIZON"
//}
