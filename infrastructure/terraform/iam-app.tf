# =============================================================================
# cluo-app IAM User
# Least-privilege credential for application/release tooling — S3 access to
# cluo-assets-prod only. Distinct from the terraform-cluo admin user.
# =============================================================================

resource "aws_iam_policy" "cluo_app" {
  name        = "cluo-app-policy"
  description = "S3 access to cluo-assets-prod only, for cluo-app"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid      = "ListAssetsProdBucket"
        Effect   = "Allow"
        Action   = "s3:ListBucket"
        Resource = aws_s3_bucket.assets_production.arn
      },
      {
        Sid    = "ReadWriteDeleteAssetsProdObjects"
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject"
        ]
        Resource = "${aws_s3_bucket.assets_production.arn}/*"
      }
    ]
  })
}

resource "aws_iam_user" "cluo_app" {
  name = "cluo-app"

  tags = {
    Name      = "cluo-app"
    Project   = var.project_name
    ManagedBy = "terraform"
  }
}

resource "aws_iam_user_policy_attachment" "cluo_app" {
  user       = aws_iam_user.cluo_app.name
  policy_arn = aws_iam_policy.cluo_app.arn
}

resource "aws_iam_access_key" "cluo_app" {
  user = aws_iam_user.cluo_app.name
}
