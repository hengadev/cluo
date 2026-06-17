# =============================================================================
# AWS CloudFront CDN
# For production assets distribution with SSL
# NOTE: Commented out for initial deployment - can be enabled later
# =============================================================================

# # =============================================================================
# # ACM Certificate for Assets Domain (must be in us-east-1 for CloudFront)
# # =============================================================================

# resource "aws_acm_certificate" "assets" {
#   provider          = aws.east
#   domain_name       = "assets.${var.domain_name}"
#   validation_method = "DNS"
#
#   subject_alternative_names = [
#     "assets-staging.${var.domain_name}"
#   ]
#
#   lifecycle {
#     create_before_destroy = true
#   }
#
#   tags = {
#     Name      = "cluo-assets-cert"
#     Project   = var.project_name
#     ManagedBy = "terraform"
#   }
# }

# # DNS validation for ACM certificate
# # Note: Validation must be completed manually by adding DNS records
# # After Terraform apply, run: terraform output acm_validation_records
# # and add those records to Cloudflare DNS

# # =============================================================================
# # CloudFront Distribution for Production Assets
# # =============================================================================

# resource "aws_cloudfront_distribution" "assets" {
#   enabled             = true
#   is_ipv6_enabled     = true
#   default_root_object = ""
#   price_class         = "PriceClass_100"  # US/Europe only (cheaper)
#
#   # Origin: S3 bucket
#   origin {
#     domain_name = aws_s3_bucket.assets_production.bucket_regional_domain_name
#     origin_id   = "S3-cluo-assets-prod"
#
#     s3_origin_config {
#       origin_access_identity = ""  # Using OAC instead
#     }
#
#     origin_access_control_id = aws_cloudfront_origin_access_control.assets.id
#   }
#
#   # Default cache behavior
#   default_cache_behavior {
#     allowed_methods  = ["GET", "HEAD", "OPTIONS"]
#     cached_methods   = ["GET", "HEAD"]
#     target_origin_id = "S3-cluo-assets-prod"
#
#     forwarded_values {
#       query_string = false
#       cookies {
#         forward = "none"
#       }
#     }
#
#     viewer_protocol_policy = "redirect-to-https"
#     min_ttl                = 0
#     default_ttl            = 86400  # 24 hours
#     max_ttl                = 31536000  # 1 year
#     compress               = true
#   }
#
#   # Custom domain
#   aliases = ["assets.${var.domain_name}"]
#
#   # Certificate for CloudFront
#   viewer_certificate {
#     acm_certificate_arn      = aws_acm_certificate.assets.arn
#     ssl_support_method       = "sni-only"
#     minimum_protocol_version = "TLSv1.2_2021"
#   }
#
#   # Custom error responses
#   custom_error_response {
#     error_code         = 403
#     response_code      = 200
#     response_page_path = "/index.html"
#   }
#
#   custom_error_response {
#     error_code         = 404
#     response_code      = 200
#     response_page_path = "/index.html"
#   }
#
#   # Restrictions (optional - for IP restriction)
#   restrictions {
#     geo_restriction {
#       restriction_type = "none"
#     }
#   }
#
#   tags = {
#     Name      = "cluo-assets-cdn"
#     Project   = var.project_name
#     ManagedBy = "terraform"
#   }
# }

# # =============================================================================
# # Origin Access Control for S3
# # =============================================================================

# resource "aws_cloudfront_origin_access_control" "assets" {
#   name                              = "cluo-assets-oac"
#   description                       = "Origin Access Control for Cluo assets bucket"
#   origin_access_control_origin_type = "s3"
#   signing_behavior                  = "always"
#   signing_protocol                  = "sigv4"
# }

# # =============================================================================
# # S3 Bucket Policy for CloudFront Access
# # =============================================================================

# resource "aws_s3_bucket_policy" "assets_production_cloudfront" {
#   bucket = aws_s3_bucket.assets_production.id
#
#   policy = jsonencode({
#     Version = "2012-10-17"
#     Statement = [
#       {
#         Sid       = "AllowCloudFrontAccess"
#         Effect    = "Allow"
#         Principal = {
#           Service = "cloudfront.amazonaws.com"
#         }
#         Action   = "s3:GetObject"
#         Resource = "${aws_s3_bucket.assets_production.arn}/*"
#         Condition = {
#           StringEquals = {
#             "AWS:SourceArn" = aws_cloudfront_distribution.assets.arn
#           }
#         }
#       }
#     ]
#   })
# }
