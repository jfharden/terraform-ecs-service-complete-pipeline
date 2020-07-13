output "ci_codebuild_project_name" {
  description = "Name of the Codebuild Project for CI builds"
  value       = aws_codebuild_project.ci.name
}

output "ci_codebuild_project_arn" {
  description = "ARN of the Codebuild Project for CI builds"
  value       = aws_codebuild_project.ci.arn
}

output "ci_codebuild_project_build_timeout" {
  description = "Build timeout (in minutes) for the CI build"
  value       = aws_codebuild_project.ci.build_timeout
}

output "ci_codebuild_project_buildspec" {
  description = "Path to the buildspec file for CI"
  value       = aws_codebuild_project.ci.buildspec
}

output "ci_codebuild_project_image" {
  description = "Path to the buildspec file for CI"
  value       = aws_codebuild_project.ci.environment.image
}

output "ecr_repository_arn" {
  description = "ARN for the ecs repository"
  value       = aws_ecr_repository.this.arn
}

output "ecr_image_scan_on_push" {
  description = "Whether scan on push is enabled for the ecr repository"
  value       = aws_ecr_repository.this.image_scanning_configuration[0].scan_on_push
}

output "ecr_image_tag_mutability" {
  description = "The tag mutability setting for the repository. One of: MUTABLE or IMMUTABLE"
  value       = aws_ecr_repository.this.image_tag_mutability
}
