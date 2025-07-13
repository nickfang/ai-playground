output "service_url" {
  description = "The URL of the service"
  value       = "http://${aws_lb.main.dns_name}"
}
