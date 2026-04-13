provider "aws" {
  region                      = "eu-central-1"
  access_key                  = "test"
  secret_key                  = "test"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  # Redirect all calls to LocalStack
  endpoints {
    dynamodb     = "http://localhost:4566"
    events       = "http://localhost:4566"
    sqs          = "http://localhost:4566"
    cloudwatch   = "http://localhost:4566"
  }
}

module "dynamodb" {
  source = "../../modules/dynamodb"
  env    = "dev"
}

module "messaging" {
  source = "../../modules/messaging"
  env    = "dev"
}