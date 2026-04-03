output "ecr_repository_url" {
  description = "ECR repository URL to push images to"
  value       = module.ecr.repository_url
}

output "ec2_public_ip" {
  description = "Public IP of the EC2 instance (set as EC2_HOST secret in GitHub)"
  value       = module.ec2.public_ip
}

output "app_url" {
  description = "URL to access the application"
  value       = "http://${module.ec2.public_ip}:8080"
}
