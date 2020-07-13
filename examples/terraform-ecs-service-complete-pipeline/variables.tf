variable "name" {
  type        = string
  description = "Name to give to all resources"
}

variable "aws_region" {
  type        = string
  description = "AWS Region"
}

variable "account_id" {
  type        = string
  description = "AWS Account ID"
}

variable "ci_build_timeout" {
  type        = number
  description = "Build timeout (in minutes) for the CI build"
  default     = 5
}

variable "ci_buildspec" {
  type        = string
  description = "Path to the buildspec file for CI"
  default     = "buildspec.ci.yaml"
}

variable "ci_codebuild_iam_policy_document_json" {
  type        = string
  description = "Additional IAM Policy document JSON to extend the default"
  default     = ""
}

variable "ci_codebuild_image" {
  type        = string
  description = "Codebuild image to use for CI build"
  default     = "aws/codebuild/amazonlinux2-x86_64-standard:3.0"
}

variable "ci_privileged_mode" {
  type        = bool
  description = "Run CI build in privileged mode"
  default     = false
}

variable "ci_compute_type" {
  type        = string
  description = "Compute type of CI build"
  default     = "BUILD_GENERAL1_SMALL"
}

variable "github_url" {
  type        = string
  description = "GitHub URL for the source"
}


variable "github_clone_depth" {
  type        = number
  description = "Depth to clone the github url"
  default     = 1
}

variable "ecr_image_scan_on_push" {
  type        = bool
  description = "Enable ECR image scanning on push"
  default     = true
}

variable "tags" {
  type        = map
  description = "Tags to apply to all resources which can be tagged"
  default     = {}
}
