# =============================================================================
# MinIO Backup Storage (staging)
# Dedicated bucket + least-privilege IAM user for the cluo-staging-backup
# container (homelab/docker/cluo/staging.docker-compose.yml), which uploads
# GPG-encrypted MinIO snapshots. Separate from cluo-assets-prod and from the
# Postgres backup buckets in backups.tf.
# =============================================================================

resource "aws_s3_bucket" "minio_backups_staging" {
  bucket = "cluo-minio-backups-staging"
}

resource "aws_s3_bucket_versioning" "minio_backups_staging" {
  bucket = aws_s3_bucket.minio_backups_staging.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "minio_backups_staging" {
  bucket = aws_s3_bucket.minio_backups_staging.id
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "minio_backups_staging" {
  bucket = aws_s3_bucket.minio_backups_staging.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_lifecycle_configuration" "minio_backups_staging" {
  bucket = aws_s3_bucket.minio_backups_staging.id

  rule {
    id     = "staging-minio-backup-retention"
    status = "Enabled"
    filter {}

    expiration {
      days = 30 # backup.sh also self-prunes after 30 days; this is a backstop
    }
  }
}

resource "aws_iam_policy" "minio_backup_staging" {
  name        = "cluo-minio-backup-staging-policy"
  description = "S3 access to cluo-minio-backups-staging only, for cluo-minio-backup-staging-user"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid      = "ListMinioBackupsStagingBucket"
        Effect   = "Allow"
        Action   = "s3:ListBucket"
        Resource = aws_s3_bucket.minio_backups_staging.arn
      },
      {
        Sid    = "ReadWriteDeleteMinioBackupsStagingObjects"
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject"
        ]
        Resource = "${aws_s3_bucket.minio_backups_staging.arn}/*"
      }
    ]
  })
}

resource "aws_iam_user" "minio_backup_staging" {
  name = "cluo-minio-backup-staging-user"

  tags = {
    Name        = "cluo-minio-backup-staging-user"
    Project     = var.project_name
    Environment = "staging"
    ManagedBy   = "terraform"
  }
}

resource "aws_iam_user_policy_attachment" "minio_backup_staging" {
  user       = aws_iam_user.minio_backup_staging.name
  policy_arn = aws_iam_policy.minio_backup_staging.arn
}

resource "aws_iam_access_key" "minio_backup_staging" {
  user = aws_iam_user.minio_backup_staging.name
}
