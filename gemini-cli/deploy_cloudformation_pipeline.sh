#!/bin/bash

# Debug prints.
# set -x

# Exit immediately if a command exits with a non-zero status.
set -e

# Define the path to the Terraform outputs JSON file
TERRAFORM_OUTPUTS_FILE="terraform_outputs.json"

# Check if the outputs file exists
if [ ! -f "$TERRAFORM_OUTPUTS_FILE" ]; then
  echo "Error: Terraform outputs file '$TERRAFORM_OUTPUTS_FILE' not found."
  echo "Please run './deploy_terraform.sh <your_github_owner>' first to generate the outputs."
  exit 1
fi

# Parse the JSON file to extract the ARNs
ECS_TASK_EXECUTION_ROLE_ARN=$(jq -r '.ecs_task_execution_role_arn.value' "$TERRAFORM_OUTPUTS_FILE")
ECS_TASK_ROLE_ARN=$(jq -r '.ecs_task_role_arn.value' "$TERRAFORM_OUTPUTS_FILE")
CONNECTION_ARN=$(jq -r '.codestar_connection_arn.value' "$TERRAFORM_OUTPUTS_FILE")

# Define CloudFormation parameters (replace YOUR_GITHUB_USERNAME and YOUR_GITHUB_REPOSITORY)
PROJECT_NAME="gemini-cli-service"
GITHUB_OWNER="nickfang" # <<< REPLACE THIS
GITHUB_REPO="ai-playground" # <<< REPLACE THIS
GITHUB_BRANCH="main"

# Get the AWS Account ID
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)

# Define the S3 bucket name based on the CloudFormation template
S3_BUCKET_NAME="$AWS_ACCOUNT_ID-codepipeline-artifact-store"

# Forcefully remove the S3 bucket and its contents
echo "Removing S3 bucket $S3_BUCKET_NAME..."
aws s3 rb "s3://$S3_BUCKET_NAME" --force || true

# Check if the stack exists before trying to delete it
if aws cloudformation describe-stacks --stack-name "$PROJECT_NAME-pipeline" >/dev/null 2>&1; then
  echo "Deleting existing CloudFormation stack '$PROJECT_NAME-pipeline'..."
  aws cloudformation delete-stack --stack-name "$PROJECT_NAME-pipeline"
  echo "Waiting for stack deletion to complete..."
  aws cloudformation wait stack-delete-complete --stack-name "$PROJECT_NAME-pipeline"
  echo "Stack deleted successfully."
else
  echo "Stack '$PROJECT_NAME-pipeline' does not exist. Skipping deletion."
fi

# Create the CloudFormation Stack for CodePipeline using the extracted ARNs
echo "Creating CloudFormation stack '$PROJECT_NAME-pipeline' with extracted ARNs..."
aws cloudformation create-stack \
  --stack-name "$PROJECT_NAME-pipeline" \
  --template-body file://.codepipeline/pipeline.yml \
  --parameters \
    ParameterKey=ProjectName,ParameterValue="$PROJECT_NAME" \
    ParameterKey=GitHubOwner,ParameterValue="$GITHUB_OWNER" \
    ParameterKey=GitHubRepo,ParameterValue="$GITHUB_REPO" \
    ParameterKey=GitHubBranch,ParameterValue="$GITHUB_BRANCH" \
    ParameterKey=ConnectionArn,ParameterValue="$CONNECTION_ARN" \
    ParameterKey=EcsTaskExecutionRoleArn,ParameterValue="$ECS_TASK_EXECUTION_ROLE_ARN" \
    ParameterKey=EcsTaskRoleArn,ParameterValue="$ECS_TASK_ROLE_ARN" \
  --capabilities CAPABILITY_IAM

echo "CloudFormation stack creation initiated. Monitor its progress in the AWS console."
