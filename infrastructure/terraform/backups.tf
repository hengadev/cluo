# =============================================================================
# PostgreSQL Backup Configuration
# Backup IAM users + policies for both environments
# =============================================================================

# =============================================================================
# Staging Backup IAM Policy
# =============================================================================

resource "aws_iam_policy" "backup_staging_policy" {
  name        = "cluo-backup-staging-policy"
  description = "IAM policy for staging PostgreSQL backups to S3"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:PutObject",
          "s3:GetObject",
          "s3:ListBucket",
          "s3:DeleteObject"
        ]
        Resource = [
          aws_s3_bucket.backups_staging.arn,
          "${aws_s3_bucket.backups_staging.arn}/*"
        ]
      }
    ]
  })
}

# =============================================================================
# Staging Backup IAM User
# =============================================================================

resource "aws_iam_user" "backup_staging_user" {
  name = "cluo-backup-staging-user"

  tags = {
    Name      = "cluo-backup-staging-user"
    Project   = var.project_name
    Environment = "staging"
    ManagedBy = "terraform"
  }
}

resource "aws_iam_user_policy_attachment" "backup_staging_attach" {
  user       = aws_iam_user.backup_staging_user.name
  policy_arn = aws_iam_policy.backup_staging_policy.arn
}

resource "aws_iam_access_key" "backup_staging_key" {
  user = aws_iam_user.backup_staging_user.name
}

# =============================================================================
# Production Backup IAM Policy
# =============================================================================

resource "aws_iam_policy" "backup_production_policy" {
  name        = "cluo-backup-production-policy"
  description = "IAM policy for production PostgreSQL backups to S3"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:PutObject",
          "s3:GetObject",
          "s3:ListBucket",
          "s3:DeleteObject"
        ]
        Resource = [
          aws_s3_bucket.backups_production.arn,
          "${aws_s3_bucket.backups_production.arn}/*"
        ]
      }
    ]
  })
}

# =============================================================================
# Production Backup IAM User
# =============================================================================

resource "aws_iam_user" "backup_production_user" {
  name = "cluo-backup-production-user"

  tags = {
    Name      = "cluo-backup-production-user"
    Project   = var.project_name
    Environment = "production"
    ManagedBy = "terraform"
  }
}

resource "aws_iam_user_policy_attachment" "backup_production_attach" {
  user       = aws_iam_user.backup_production_user.name
  policy_arn = aws_iam_policy.backup_production_policy.arn
}

resource "aws_iam_access_key" "backup_production_key" {
  user = aws_iam_user.backup_production_user.name
}

# =============================================================================
# Assets IAM Policy for Application Uploads
# =============================================================================

resource "aws_iam_policy" "assets_staging_policy" {
  name        = "cluo-assets-staging-policy"
  description = "IAM policy for staging assets S3 access"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:PutObject",
          "s3:GetObject",
          "s3:DeleteObject"
        ]
        Resource = "${aws_s3_bucket.assets_staging.arn}/*"
      }
    ]
  })
}

resource "aws_iam_user" "assets_staging_user" {
  name = "cluo-assets-staging-user"

  tags = {
    Name      = "cluo-assets-staging-user"
    Project   = var.project_name
    Environment = "staging"
    ManagedBy = "terraform"
  }
}

resource "aws_iam_user_policy_attachment" "assets_staging_attach" {
  user       = aws_iam_user.assets_staging_user.name
  policy_arn = aws_iam_policy.assets_staging_policy.arn
}

resource "aws_iam_access_key" "assets_staging_key" {
  user = aws_iam_user.assets_staging_user.name
}

# =============================================================================
# Production Assets IAM Policy
# =============================================================================

resource "aws_iam_policy" "assets_production_policy" {
  name        = "cluo-assets-production-policy"
  description = "IAM policy for production assets S3 access"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:PutObject",
          "s3:GetObject",
          "s3:DeleteObject"
        ]
        Resource = "${aws_s3_bucket.assets_production.arn}/*"
      }
    ]
  })
}

resource "aws_iam_user" "assets_production_user" {
  name = "cluo-assets-production-user"

  tags = {
    Name      = "cluo-assets-production-user"
    Project   = var.project_name
    Environment = "production"
    ManagedBy = "terraform"
  }
}

resource "aws_iam_user_policy_attachment" "assets_production_attach" {
  user       = aws_iam_user.assets_production_user.name
  policy_arn = aws_iam_policy.assets_production_policy.arn
}

resource "aws_iam_access_key" "assets_production_key" {
  user = aws_iam_user.assets_production_user.name
}
