# =============================================================================
# HashiCorp Vault Configuration
# Auto-unseal using AWS KMS + S3 storage backend
# =============================================================================

# =============================================================================
# KMS Key for Vault Auto-Unseal
# =============================================================================

resource "aws_kms_key" "vault" {
  description             = "Cluo Vault auto-unseal key"
  deletion_window_in_days = 30
  enable_key_rotation     = true

  tags = {
    Name      = "cluo-vault-unseal"
    Project   = var.project_name
    ManagedBy = "terraform"
  }
}

resource "aws_kms_alias" "vault" {
  name          = "alias/cluo-vault-unseal"
  target_key_id = aws_kms_key.vault.key_id
}

# =============================================================================
# IAM Policy for Vault S3 + KMS Access
# =============================================================================

resource "aws_iam_policy" "vault_policy" {
  name        = "cluo-vault-policy"
  description = "IAM policy for Vault S3 storage and KMS unseal"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:ListBucket",
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject"
        ]
        Resource = [
          aws_s3_bucket.vault.arn,
          "${aws_s3_bucket.vault.arn}/*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "kms:Encrypt",
          "kms:Decrypt",
          "kms:DescribeKey"
        ]
        Resource = [
          aws_kms_key.vault.arn
        ]
      }
    ]
  })
}

# =============================================================================
# IAM User for Vault
# =============================================================================

resource "aws_iam_user" "vault_user" {
  name = "cluo-vault-user"

  tags = {
    Name      = "cluo-vault-user"
    Project   = var.project_name
    ManagedBy = "terraform"
  }
}

resource "aws_iam_user_policy_attachment" "vault_attach" {
  user       = aws_iam_user.vault_user.name
  policy_arn = aws_iam_policy.vault_policy.arn
}

# =============================================================================
# IAM Access Key for Vault
# =============================================================================

resource "aws_iam_access_key" "vault_key" {
  user = aws_iam_user.vault_user.name
}
