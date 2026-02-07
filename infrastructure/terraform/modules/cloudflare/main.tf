# API subdomain
resource "cloudflare_record" "api" {
  zone_id = var.zone_id
  name    = "api"
  value   = var.server_ip
  type    = "A"
  ttl     = 1
  proxied = var.proxy_api
}

# Web subdomain
resource "cloudflare_record" "web" {
  zone_id = var.zone_id
  name    = "app"
  value   = var.server_ip
  type    = "A"
  ttl     = 1
  proxied = var.proxy_web
}

# Mobile subdomain
resource "cloudflare_record" "mobile" {
  zone_id = var.zone_id
  name    = "mobile"
  value   = var.server_ip
  type    = "A"
  ttl     = 1
  proxied = var.proxy_mobile
}

# Wildcard for additional subdomains
resource "cloudflare_record" "wildcard" {
  zone_id = var.zone_id
  name    = "*"
  value   = var.server_ip
  type    = "A"
  ttl     = 1
  proxied = true
}
