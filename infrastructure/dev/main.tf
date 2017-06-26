variable "project" {
  type = "string"
  default = "standup"
}

variable "aws_region" {
  type = "string"
  default = "us-east-1"
}

variable "apex_environment" {
  type = "string"
  default = "dev"
}

#
# private subnet ranges
#
variable "private_cidr" {
  type = "list"
  description = "private subnet CIDR ranges for VPC"
  default = ["10.0.0.0/20", "10.0.16.0/20", "10.0.32.0/20"]
}

#
# IAM users and roles
#
module "iam" {
  source = "../modules/iam"
  project = "${var.project}"
  environment = "${var.apex_environment}"
}

#
# VPC
# see http://www.vlsm-calc.net/ for size calculations
#
module "vpc" {
  source = "github.com/terraform-community-modules/tf_aws_vpc"

  name                 = "${var.project}-${var.apex_environment}-vpc"
  azs                  = ["us-east-1a", "us-east-1b", "us-east-1c"]
  cidr                 = "10.0.0.0/16" # 65534 IP addresses
  public_subnets       = ["10.0.96.0/20", "10.0.112.0/20", "10.0.128.0/20"]
  private_subnets      = ["${var.private_cidr}"]
  enable_nat_gateway   = "true"
  enable_dns_support   = "true"
  enable_dns_hostnames = "true"
  tags {
    "Terraform"        = "true"
    "Environment"      = "${var.apex_environment}"
    "Project"          = "${var.project}"
  }
}
