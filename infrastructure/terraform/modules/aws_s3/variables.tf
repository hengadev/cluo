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

variable "cors_allowed_origins" {
  description = "List of allowed origins for CORS (e.g., [\"https://app.example.com\"]). Empty list disables CORS."
  type        = list(string)
  default     = []
}

variable "lifecycle_rules" {
  description = "Custom lifecycle rules. When empty, defaults to deleting old versions after 30 days."
  type = list(object({
    id                                 = string
    status                             = string
    noncurrent_version_expiration_days = optional(number)
    transition_days                    = optional(number)
    transition_storage_class           = optional(string)
    expiration_days                    = optional(number)
  }))
  default = []
}
