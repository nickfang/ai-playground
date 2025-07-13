output "service_url" {
  description = "The URL of the service"
  value       = "http://${aws_lb.main.dns_name}"
}

output "codestar_connection_arn" {
  description = "The ARN of the CodeStar Connection to GitHub"
  value       = aws_codestarconnections_connection.github.arn
}