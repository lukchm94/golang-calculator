resource "aws_dynamodb_table" "calculations" {
  name           = "${var.env}-calculations"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "PartitionKey" # e.g., CALC#<ID>
  range_key      = "SortKey" # e.g., METADATA

  attribute {
    name = "PartitionKey"
    type = "S"
  }

  attribute {
    name = "SortKey"
    type = "S"
  }

  tags = {
    Environment = var.env
  }
}

variable "env" { type = string }