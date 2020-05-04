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
  value = aws_ecr_repository.this.image_tag_mutability
}

output "ecr_tags" {
  description = "Tags applied to the ECR repository"
  value = aws_ecr_repository.this.tags
}
