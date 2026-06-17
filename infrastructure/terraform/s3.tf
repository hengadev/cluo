# =============================================================================
# AWS S3 Resources
# Buckets for: Assets (staging + production), Backups, Vault
# =============================================================================

# =============================================================================
# Staging Assets Bucket
# =============================================================================

resource "aws_s3_bucket" "assets_staging" {
  bucket = "cluo-assets-staging"
}

resource "aws_s3_bucket_versioning" "assets_staging" {
  bucket = aws_s3_bucket.assets_staging.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_cors_configuration" "assets_staging" {
  bucket = aws_s3_bucket.assets_staging.id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "HEAD"]
    allowed_origins = [
      "http://localhost:5173",
      "https://staging.${var.domain_name}"
    ]
    expose_headers  = ["ETag"]
    max_age_seconds = 3600
  }
}

resource "aws_s3_bucket_ownership_controls" "assets_staging" {
  bucket = aws_s3_bucket.assets_staging.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_public_access_block" "assets_staging" {
  bucket = aws_s3_bucket.assets_staging.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_acl" "assets_staging" {
  depends_on = [aws_s3_bucket_ownership_controls.assets_staging]
  bucket     = aws_s3_bucket.assets_staging.id
  acl        = "public-read"
}

# =============================================================================
# Production Assets Bucket (with CloudFront)
# =============================================================================

resource "aws_s3_bucket" "assets_production" {
  bucket = "cluo-assets-prod"
}

resource "aws_s3_bucket_versioning" "assets_production" {
  bucket = aws_s3_bucket.assets_production.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_cors_configuration" "assets_production" {
  bucket = aws_s3_bucket.assets_production.id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "HEAD"]
    allowed_origins = [
      "https://${var.domain_name}",
      "https://www.${var.domain_name}",
      "https://mobile.${var.domain_name}"
    ]
    expose_headers  = ["ETag"]
    max_age_seconds = 3600
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "assets_production" {
  bucket = aws_s3_bucket.assets_production.id
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_ownership_controls" "assets_production" {
  bucket = aws_s3_bucket.assets_production.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_public_access_block" "assets_production" {
  bucket = aws_s3_bucket.assets_production.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_acl" "assets_production" {
  depends_on = [aws_s3_bucket_ownership_controls.assets_production]
  bucket     = aws_s3_bucket.assets_production.id
  acl        = "public-read"
}

# =============================================================================
# Staging PostgreSQL Backup Bucket
# =============================================================================

resource "aws_s3_bucket" "backups_staging" {
  bucket = "cluo-postgres-backups-staging"
}

resource "aws_s3_bucket_versioning" "backups_staging" {
  bucket = aws_s3_bucket.backups_staging.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "backups_staging" {
  bucket = aws_s3_bucket.backups_staging.id
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "backups_staging" {
  bucket = aws_s3_bucket.backups_staging.id

  rule {
    id      = "staging-backup-retention"
    status  = "Enabled"
    filter {}

    expiration {
      days = 30  # Staging backups retained for 30 days
    }
  }
}

# =============================================================================
# Production PostgreSQL Backup Bucket
# =============================================================================

resource "aws_s3_bucket" "backups_production" {
  bucket = "cluo-postgres-backups-production"
}

resource "aws_s3_bucket_versioning" "backups_production" {
  bucket = aws_s3_bucket.backups_production.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "backups_production" {
  bucket = aws_s3_bucket.backups_production.id
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "backups_production" {
  bucket = aws_s3_bucket.backups_production.id

  rule {
    id      = "production-backup-retention"
    status  = "Enabled"
    filter {}

    noncurrent_version_expiration {
      noncurrent_days = 90  # Daily backups retained for 90 days
    }
  }
}

# =============================================================================
# Vault Storage Bucket
# =============================================================================

resource "aws_s3_bucket" "vault" {
  bucket = "cluo-vault-storage"
}

resource "aws_s3_bucket_versioning" "vault" {
  bucket = aws_s3_bucket.vault.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "vault" {
  bucket = aws_s3_bucket.vault.id
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "vault" {
  bucket = aws_s3_bucket.vault.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}
