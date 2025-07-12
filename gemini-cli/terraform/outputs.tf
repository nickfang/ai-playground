output "service_url" {
  description = "The URL of the service"
  value       = "http://${aws_ecs_service.main.network_configuration[0].assign_public_ip ? aws_ecs_service.main.network_configuration[0].public_ip : ""}:3000"
}
