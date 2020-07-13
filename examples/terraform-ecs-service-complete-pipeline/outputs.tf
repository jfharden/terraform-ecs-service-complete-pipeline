output "ci_codebuild_project_name" {
  description = "Name of the Codebuild Project for CI builds"
  value       = module.service_pipeline.ci_codebuild_project_name
}

output "ci_codebuild_project_arn" {
  description = "ARN of the Codebuild Project for CI builds"
  value       = module.service_pipeline.ci_codebuild_project_arn
}

output "ci_build_timeout" {
  description = "Build timeout (in minutes) for the CI build"
  value       = module.service_pipeline.ci_codebuild_project_build_timeout
}

output "ci_buildspec" {
  description = "Path to the buildspec file for CI"
  value       = module.service_pipeline.ci_codebuild_project_buildspec
}

output "ci_codebuild_image" {
  description = "Path to the buildspec file for CI"
  value       = module.service_pipeline.ci_codebuild_project_image
}

output "ecr_repository_arn" {
  value = module.service_pipeline.ecr_repository_arn
}

output "ecr_image_scan_on_push" {
  value = module.service_pipeline.ecr_image_scan_on_push
}

output "ecr_image_tag_mutability" {
  description = "The tag mutability setting for the repository. One of: MUTABLE or IMMUTABLE"
  value = module.service_pipeline.ecr_image_tag_mutability
}
