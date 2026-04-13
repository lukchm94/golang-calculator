# The Dead Letter Queue
resource "aws_sqs_queue" "dlq" {
  name = "${var.env}-calculator-dlq"
}

# DLQ for Delivery/Forwarding Failures (EventBridge level)
resource "aws_sqs_queue" "eb_delivery_dlq" {
  name = "${var.env}-calculator-delivery-dlq"
}

# The Main SQS Queue
resource "aws_sqs_queue" "main_queue" {
  name = "${var.env}-calculator-main-queue"
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.dlq.arn
    maxReceiveCount     = 4
  })
}

# The Custom Event Bus
resource "aws_cloudwatch_event_bus" "bus" {
  name = "${var.env}-calculator-event-bus"
}

# Rule to catch specific events
resource "aws_cloudwatch_event_rule" "rule" {
  name           = "${var.env}-calc-to-sqs-rule"
  event_bus_name = aws_cloudwatch_event_bus.bus.name
  event_pattern  = jsonencode({
    source = ["calculator.app"]
  })
}

# Link Rule -> SQS
resource "aws_cloudwatch_event_target" "target" {
  rule           = aws_cloudwatch_event_rule.rule.name
  event_bus_name = aws_cloudwatch_event_bus.bus.name
  arn            = aws_sqs_queue.main_queue.arn

  dead_letter_config {
    arn = aws_sqs_queue.eb_delivery_dlq.arn
  }

  retry_policy {
    maximum_event_age_in_seconds = 86400 # 24 hours
    maximum_retry_attempts       = 4
  }
}

# Define variables.tf in the same folder
variable "env" { type = string }