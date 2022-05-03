terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.12.1"
    }
  }
}

locals {
  current_region = "us-east-2" #US East (Ohio)
  key_pair_name  = "terraform-lab-key-pair"
}

provider "aws" {
  region     = local.current_region
  access_key = "access_key" #terraform-lab
  secret_key = "secret_key"
  # shared_credentials_file = "$HOME/.aws/credentials"
  # profile                 = "1edge-1edge-creator-01"
}

variable "default_tags_list" {
  type = list(object({
    key                 = string
    value               = string
    propagate_at_launch = bool
  }))
  default = [
    {
      key                 = "ApplicationName"
      value               = "NSBU Terraform Session"
      propagate_at_launch = true
    },
    {
      key                 = "EnvironmentName"
      value               = "NSBU-Terraform-Session-Lab"
      propagate_at_launch = true
    },
    {
      key                 = "CiscoMailAlias"
      value               = "lpereir2@cisco.com"
      propagate_at_launch = true
    },
    {
      key                 = "ResourceOwner"
      value               = "NSBU-Lab-Session"
      propagate_at_launch = true
    }
  ]
}

locals {
  default_tags = {
    for obj in var.default_tags_list : "${obj.key}" => obj.value
  }
}

resource "aws_vpc" "terraform_lab_vpc" {
  cidr_block           = "10.0.0.0/16"
  instance_tenancy     = "default"
  enable_dns_support   = true
  enable_dns_hostnames = true
  tags = merge(
    local.default_tags,
    {
      Name = "terraform_lab_vpc"
    }
  )
}

resource "aws_subnet" "terraform_lab_subnet_01" {
  vpc_id            = aws_vpc.terraform_lab_vpc.id
  cidr_block        = "10.0.0.0/24"
  availability_zone = "${local.current_region}a"

  tags = {
    "Name" = "terraform_lab_subnet_01"
  }
}

resource "aws_subnet" "terraform_lab_subnet_02" {
  vpc_id            = aws_vpc.terraform_lab_vpc.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = "${local.current_region}b"
  tags = {
    Name = "terraform_lab_subnet_02"
  }
}

resource "aws_internet_gateway" "terraform_lab_internet_gateway" {
  vpc_id = aws_vpc.terraform_lab_vpc.id
  tags = {
    Name = "terraform_lab_internet_gateway"
  }
}

resource "aws_route" "terraform_lab_internet_gateway_route" {
  route_table_id         = aws_vpc.terraform_lab_vpc.default_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.terraform_lab_internet_gateway.id
}

resource "aws_security_group" "terraform_lab_security_group" {
  name        = "terraform_lab_security_group"
  description = "SG for Terraform Lab EC2"
  vpc_id      = aws_vpc.terraform_lab_vpc.id

  ingress {
    description = "All for http"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "terraform_lab_security_group"
  }
}

data "aws_ami" "amazon_linux_2" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-kernel-5.10-hvm-2.0.20220419.0-x86_64-gp2"]
  }
}

resource "aws_key_pair" "terraform_lab_key_pair" {
  key_name   = local.key_pair_name
  public_key = "ssh-rsa ......= lpereir2@LPEREIR2-M-J62P"
}

resource "aws_instance" "terraform_lab_ec2" {
  ami                         = data.aws_ami.amazon_linux_2.id
  instance_type               = "t2.micro"
  key_name                    = local.key_pair_name
  associate_public_ip_address = "true"
  vpc_security_group_ids      = [aws_security_group.terraform_lab_security_group.id]
  subnet_id                   = aws_subnet.terraform_lab_subnet_01.id
  root_block_device {
    volume_type           = "gp2"
    volume_size           = 8
    delete_on_termination = true
  }

  user_data = <<-EOF
              #!/bin/bash
              yum update -y
              yum install -y httpd
              systemctl start httpd
              systemctl enable httpd
              echo "<h1>Hello world from $(hostname -f)</h1>" > /var/www/html/index.html
              EOF

  tags = merge(
    local.default_tags,
    {
      Name = "Terraform Lab EC2"
    }
  )
}

output "EC2_public_dns" {
  value = aws_instance.terraform_lab_ec2.public_dns
}

output "EC2_private_dns" {
  value = aws_instance.terraform_lab_ec2.private_dns
}
