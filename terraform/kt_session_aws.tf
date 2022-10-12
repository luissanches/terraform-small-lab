terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~>4.34.0"
    }

    null = {
      source  = "hashicorp/null"
      version = "3.1.1"
    }
  }
  required_version = ">=1.3.1"
}

locals {
  region_ohio   = "us-east-2"
  region_oregon = "us-west-2"
  key_pair_name = "terraform-lab-key-pair"
}

provider "null" {
  # Configuration options
}

provider "aws" {
  region = local.region_ohio
}


resource "null_resource" "echo" {
	count = 0
  provisioner "local-exec" {
    command = "echo hello..."
  }
}

