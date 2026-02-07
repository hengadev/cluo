output "bucket_name" {
  description = "Name of the S3 bucket"
  value       = aws_s3_bucket.media.id
}

output "bucket_arn" {
  description = "ARN of the S3 bucket"
  value       = aws_s3_bucket.media.arn
}

output "bucket_region" {
  description = "Region of the S3 bucket"
  value       = var.region
}

output "iam_access_key_id" {
  description = "IAM access key ID"
  value       = var.create_iam_user ? aws_iam_access_key.app[0].id : null
  sensitive   = true
}

output "iam_secret_access_key" {
  description = "IAM secret access key"
  value       = var.create_iam_user ? aws_iam_access_key.app[0].secret : null
  sensitive   = true
}

output "iam_user_name" {
  description = "IAM user name"
  value       = var.create_iam_user ? aws_iam_user.app[0].name : null
}
