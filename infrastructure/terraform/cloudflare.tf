# =============================================================================
# Cloudflare DNS Records
# Single VPS Architecture: Staging + Production
# =============================================================================

locals {
  vps_ipv4 = hcloud_server.main.ipv4_address
}

# =============================================================================
# Staging Environment DNS Records
# =============================================================================

resource "cloudflare_dns_record" "staging_web" {
  zone_id = var.zone_id
  name    = "staging"
  type    = "A"
  content = local.vps_ipv4
  proxied = true
  ttl     = 1
}

resource "cloudflare_dns_record" "staging_api" {
  zone_id = var.zone_id
  name    = "staging-api"
  type    = "A"
  content = local.vps_ipv4
  proxied = true
  ttl     = 1
}

resource "cloudflare_dns_record" "staging_mobile" {
  zone_id = var.zone_id
  name    = "staging-mobile"
  type    = "A"
  content = local.vps_ipv4
  proxied = true
  ttl     = 1
}

resource "cloudflare_dns_record" "staging_assets" {
  zone_id = var.zone_id
  name    = "assets-staging"
  type    = "CNAME"
  content = "cluo-assets-staging.s3.amazonaws.com"
  proxied = false
  ttl     = 3600
}

# =============================================================================
# Production Environment DNS Records
# =============================================================================

resource "cloudflare_dns_record" "production_web" {
  zone_id = var.zone_id
  name    = var.domain_name
  type    = "A"
  content = local.vps_ipv4
  proxied = true
  ttl     = 1
}

resource "cloudflare_dns_record" "production_www" {
  zone_id = var.zone_id
  name    = "www"
  type    = "CNAME"
  content = var.domain_name
  proxied = true
  ttl     = 1
}

resource "cloudflare_dns_record" "production_api" {
  zone_id = var.zone_id
  name    = "api"
  type    = "A"
  content = local.vps_ipv4
  proxied = true
  ttl     = 1
}

resource "cloudflare_dns_record" "production_mobile" {
  zone_id = var.zone_id
  name    = "mobile"
  type    = "A"
  content = local.vps_ipv4
  proxied = true
  ttl     = 1
}

resource "cloudflare_dns_record" "production_assets" {
  zone_id = var.zone_id
  name    = "assets"
  type    = "CNAME"
  content = aws_cloudfront_distribution.assets.domain_name
  proxied = false
  ttl     = 3600
}

# =============================================================================
# Email DNS Records (MX, SPF, DMARC)
# =============================================================================

# MX Records
resource "cloudflare_dns_record" "mx" {
  for_each = {
    for idx, mx in var.mx_servers : idx => mx
  }

  zone_id  = var.zone_id
  name     = var.domain_name
  type     = "MX"
  content  = each.value.server
  priority = each.value.priority
  proxied  = false
}

# SPF Record
resource "cloudflare_dns_record" "spf" {
  zone_id = var.zone_id
  name    = var.domain_name
  type    = "TXT"
  content = "\"v=spf1 include:amazonses.com ~all\""
  proxied = false
  ttl     = 3600
}

# DMARC Record
resource "cloudflare_dns_record" "dmarc" {
  zone_id = var.zone_id
  name    = "_dmarc"
  type    = "TXT"
  content = "\"v=DMARC1; p=${var.email_dmarc_policy}; rua=mailto:${var.contact_email}\""
  proxied = false
  ttl     = 3600
}

# =============================================================================
# SES Email Verification Records
# =============================================================================

# SES Domain Verification
resource "cloudflare_dns_record" "ses_verification" {
  zone_id = var.zone_id
  name    = "_amazonses.${var.domain_name}"
  type    = "TXT"
  content = "\"${aws_ses_domain_identity.main.verification_token}\""
  proxied = false
  ttl     = 3600
}

# SES DKIM Records (3 tokens)
resource "cloudflare_dns_record" "ses_dkim" {
  count = 3

  zone_id = var.zone_id
  name    = "${aws_ses_domain_dkim.main.dkim_tokens[count.index]}._domainkey"
  type    = "CNAME"
  content = "${aws_ses_domain_dkim.main.dkim_tokens[count.index]}.dkim.amazonses.com"
  proxied = false
  ttl     = 3600
}

# Mailbox Provider DKIM Records
resource "cloudflare_dns_record" "mailbox_dkim" {
  for_each = var.email_dkim_records

  zone_id = var.zone_id
  name    = each.key
  type    = each.value.type
  content = each.value.type == "TXT" ? "\"${each.value.content}\"" : each.value.content
  proxied = false
  ttl     = 3600
}

# =============================================================================
# ACM Certificate Validation for CloudFront
# =============================================================================

resource "cloudflare_dns_record" "acm_validation" {
  for_each = {
    for dvo in aws_acm_certificate_validation.assets : dvo.domain_validation_options[0].resource_record_name => dvo
  }

  zone_id = var.zone_id
  name    = each.value.domain_validation_options[0].resource_record_name
  type    = each.value.domain_validation_options[0].resource_record_type
  content = each.value.domain_validation_options[0].resource_record_value
  proxied = false
  ttl     = 60
}
