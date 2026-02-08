# Hetzner outputs
output "server_ip" {
  description = "Public IP address of the Hetzner server"
  value       = module.hetzner.server_ip
}

output "server_name" {
  description = "Name of the Hetzner server"
  value       = module.hetzner.server_name
}

output "server_ssh_command" {
  description = "SSH command to connect to the server"
  value       = "root@${module.hetzner.server_ip}"
}

# Cloudflare outputs
output "dns_records" {
  description = "DNS records created in Cloudflare"
  value       = module.cloudflare.dns_records
}

# AWS S3 media bucket outputs
output "s3_bucket_name" {
  description = "Name of the S3 media bucket"
  value       = module.aws_s3.bucket_name
}

output "s3_bucket_arn" {
  description = "ARN of the S3 media bucket"
  value       = module.aws_s3.bucket_arn
}

output "iam_access_key_id" {
  description = "IAM access key ID for S3 media access"
  value       = module.aws_s3.iam_access_key_id
  sensitive   = true
}

output "iam_secret_access_key" {
  description = "IAM secret access key for S3 media access"
  value       = module.aws_s3.iam_secret_access_key
  sensitive   = true
}

# AWS S3 backup bucket outputs
output "backup_bucket_name" {
  description = "Name of the S3 backup bucket"
  value       = module.aws_s3_backup.bucket_name
}

output "backup_bucket_arn" {
  description = "ARN of the S3 backup bucket"
  value       = module.aws_s3_backup.bucket_arn
}

output "backup_iam_access_key_id" {
  description = "IAM access key ID for S3 backup access"
  value       = module.aws_s3_backup.iam_access_key_id
  sensitive   = true
}

output "backup_iam_secret_access_key" {
  description = "IAM secret access key for S3 backup access"
  value       = module.aws_s3_backup.iam_secret_access_key
  sensitive   = true
}

# Post-provisioning instructions
output "next_steps" {
  description = "Next steps after Terraform apply"
  value       = <<-EOT
    1. Connect to the server: ssh root@${module.hetzner.server_ip}
    2. Install Docker and Docker Compose
    3. Deploy your application using docker-compose.yml
    4. Configure your .env with the S3 credentials from terraform output
  EOT
}
