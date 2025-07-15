#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Define the path to the Terraform directory
TERRAFORM_DIR="terraform"

# Define the output file for Terraform ARNs
TERRAFORM_OUTPUTS_FILE="../terraform_outputs.json"

# Check for GitHub owner variable
if [ -z "$1" ]; then
  echo "Usage: $0 <github_owner>"
  echo "Please provide your GitHub owner/organization name as the first argument."
  exit 1
fi

GITHUB_OWNER="$1"

echo "Navigating to Terraform directory: $TERRAFORM_DIR"
cd "$TERRAFORM_DIR"

echo "Initializing Terraform..."
terraform init

echo "Applying Terraform changes..."
# The -auto-approve flag is used for automation. For interactive use, remove it.
# If you want to review the plan before applying, run 'terraform plan' first.
terraform apply -var="github_owner=$GITHUB_OWNER"

echo "Generating Terraform outputs to $TERRAFORM_OUTPUTS_FILE..."
terraform output -json > "$TERRAFORM_OUTPUTS_FILE"

echo "Returning to project root..."
cd ..

echo "Terraform deployment complete. Outputs saved to $TERRAFORM_OUTPUTS_FILE."
