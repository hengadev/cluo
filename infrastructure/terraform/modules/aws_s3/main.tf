# S3 Bucket for media storage
resource "aws_s3_bucket" "media" {
  bucket = var.bucket_name

  tags = {
    Name = "${var.bucket_name}-media"
  }
}

# Bucket versioning (optional)
resource "aws_s3_bucket_versioning" "media" {
  bucket = aws_s3_bucket.media.id

  versioning_configuration {
    status = "Enabled"
  }
}

# Bucket server-side encryption
resource "aws_s3_bucket_server_side_encryption_configuration" "media" {
  bucket = aws_s3_bucket.media.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# Bucket public access block (disable if you need public access)
resource "aws_s3_bucket_public_access_block" "media" {
  bucket = aws_s3_bucket.media.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# Bucket lifecycle policy
resource "aws_s3_bucket_lifecycle_configuration" "media" {
  bucket = aws_s3_bucket.media.id

  # Default rule: delete old versions after 30 days
  dynamic "rule" {
    for_each = length(var.lifecycle_rules) == 0 ? [1] : []
    content {
      id     = "delete-old-versions"
      status = "Enabled"

      noncurrent_version_expiration {
        noncurrent_days = 30
      }
    }
  }

  # Custom lifecycle rules
  dynamic "rule" {
    for_each = var.lifecycle_rules
    content {
      id     = rule.value.id
      status = rule.value.status

      dynamic "noncurrent_version_expiration" {
        for_each = rule.value.noncurrent_version_expiration_days != null ? [1] : []
        content {
          noncurrent_days = rule.value.noncurrent_version_expiration_days
        }
      }

      dynamic "transition" {
        for_each = rule.value.transition_days != null ? [1] : []
        content {
          days          = rule.value.transition_days
          storage_class = rule.value.transition_storage_class
        }
      }

      dynamic "expiration" {
        for_each = rule.value.expiration_days != null ? [1] : []
        content {
          days = rule.value.expiration_days
        }
      }
    }
  }
}

# CORS configuration for presigned GET URLs
resource "aws_s3_bucket_cors_configuration" "media" {
  count  = length(var.cors_allowed_origins) > 0 ? 1 : 0
  bucket = aws_s3_bucket.media.id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "HEAD"]
    allowed_origins = var.cors_allowed_origins
    max_age_seconds = 3600
  }
}

# IAM user for application access
resource "aws_iam_user" "app" {
  count = var.create_iam_user ? 1 : 0
  name  = var.iam_user_name

  tags = {
    Description = "IAM user for CLUO application S3 access"
  }
}

# IAM access key for the user
resource "aws_iam_access_key" "app" {
  count = var.create_iam_user ? 1 : 0
  user  = aws_iam_user.app[0].name
}

# IAM policy for S3 access
resource "aws_iam_user_policy" "s3_access" {
  count = var.create_iam_user ? 1 : 0
  name  = "${var.bucket_name}-s3-policy"
  user  = aws_iam_user.app[0].name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:PutObject",
          "s3:GetObject",
          "s3:DeleteObject",
          "s3:ListBucket"
        ]
        Resource = [
          aws_s3_bucket.media.arn,
          "${aws_s3_bucket.media.arn}/*"
        ]
      }
    ]
  })
}
