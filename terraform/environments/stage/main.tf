provider "aws" {
  region = "eu-central-1"
  # Authentication usually handled via IAM Roles or ENV vars in CI/CD
}

module "database" {
  source = "../../modules/dynamodb"
  env    = "stage"
}

module "messaging" {
  source = "../../modules/messaging"
  env    = "stage"
}