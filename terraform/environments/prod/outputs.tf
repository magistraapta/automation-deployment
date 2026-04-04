output "ecr_repository_url" {
  value = module.ecr.repository_url
}

output "ec2_public_ip" {
  value = module.ec2.public_ip
}

output "app_url" {
  value = "http://${module.ec2.public_ip}:8080"
}

output "rds_endpoint" {
  value = module.rds.endpoint
}
