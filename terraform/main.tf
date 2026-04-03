terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
  profile = "devops-automation"
}

module "ecr" {
  source          = "./modules/ecr"
  repository_name = var.ecr_repository_name
}

module "networking" {
  source       = "./modules/networking"
  project_name = var.project_name
}

module "ec2" {
  source            = "./modules/ec2"
  project_name      = var.project_name
  vpc_id            = module.networking.vpc_id
  subnet_id         = module.networking.public_subnet_id
  security_group_id = module.networking.security_group_id
  ecr_registry_url  = module.ecr.repository_url
  key_name          = var.key_name
  instance_type     = var.instance_type
}
