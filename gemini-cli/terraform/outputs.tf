output "service_url" {
  description = "The URL of the service"
  value       = "http://${aws_lb.main.dns_name}"
}

output "codestar_connection_arn" {
  description = "The ARN of the CodeStar Connection to GitHub"
  value       = aws_codestarconnections_connection.github.arn
}

output "ecs_task_execution_role_arn" {
  description = "The ARN of the ECS Task Execution Role"
  value       = aws_iam_role.ecs_task_execution.arn
}

output "ecs_task_role_arn" {
  description = "The ARN of the ECS Task Role"
  value       = aws_iam_role.ecs_task.arn
}