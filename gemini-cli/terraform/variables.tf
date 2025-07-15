variable "aws_region" {
  description = "The AWS region to deploy the infrastructure"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "The name of the project"
  type        = string
  default     = "gemini-cli-service"
}

variable "github_owner" {
  description = "The GitHub owner or organization for the CodeStar Connection"
  type        = string
}


