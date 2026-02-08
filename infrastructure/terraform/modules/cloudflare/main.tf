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

# Zone-level SSL/TLS and security settings
resource "cloudflare_zone_settings_override" "settings" {
  zone_id = var.zone_id

  settings {
    ssl                      = "strict"
    min_tls_version          = "1.2"
    always_use_https         = "on"
    automatic_https_rewrites = "on"

    security_header {
      enabled            = true
      max_age            = 31536000
      include_subdomains = true
      nosniff            = true
    }
  }
}

# Cache rules: bypass API, cache static assets on app/mobile
resource "cloudflare_ruleset" "cache_rules" {
  count = var.enable_cache_rules ? 1 : 0

  zone_id     = var.zone_id
  name        = "Cache rules"
  description = "Bypass cache for API, cache static assets"
  kind        = "zone"
  phase       = "http_request_cache_settings"

  # Bypass cache for API subdomain
  rules {
    action = "set_cache_settings"
    action_parameters {
      cache = false
    }
    expression  = "(http.host eq \"api.${var.domain_name}\")"
    description = "Bypass cache for API"
    enabled     = true
  }

  # Cache static assets on app and mobile subdomains
  rules {
    action = "set_cache_settings"
    action_parameters {
      cache = true
      edge_ttl {
        mode    = "override_origin"
        default = 86400
      }
    }
    expression  = "(http.host eq \"app.${var.domain_name}\" or http.host eq \"mobile.${var.domain_name}\") and http.request.uri.path.extension in {\"js\" \"css\" \"svg\" \"png\" \"jpg\" \"jpeg\" \"webp\" \"gif\" \"ico\" \"woff\" \"woff2\" \"ttf\" \"eot\"}"
    description = "Cache static assets on app and mobile"
    enabled     = true
  }
}
