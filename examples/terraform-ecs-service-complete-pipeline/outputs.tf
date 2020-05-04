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

output "ecr_tags" {
  value = module.service_pipeline.ecr_tags
}
