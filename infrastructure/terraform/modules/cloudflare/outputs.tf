output "dns_records" {
  description = "DNS records created"
  value = {
    api    = cloudflare_record.api.hostname
    web    = cloudflare_record.web.hostname
    mobile = cloudflare_record.mobile.hostname
    wildcard = cloudflare_record.wildcard.hostname
  }
}

output "api_hostname" {
  description = "API hostname"
  value       = cloudflare_record.api.hostname
}

output "web_hostname" {
  description = "Web hostname"
  value       = cloudflare_record.web.hostname
}

output "mobile_hostname" {
  description = "Mobile hostname"
  value       = cloudflare_record.mobile.hostname
}
