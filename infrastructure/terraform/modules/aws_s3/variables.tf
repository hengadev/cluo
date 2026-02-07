variable "bucket_name" {
  description = "Name of the S3 bucket"
  type        = string
}

variable "region" {
  description = "AWS region"
  type        = string
}

variable "create_iam_user" {
  description = "Whether to create an IAM user for S3 access"
  type        = bool
  default     = true
}

variable "iam_user_name" {
  description = "Name of the IAM user to create"
  type        = string
  default     = "cluo-app"
}
