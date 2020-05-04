resource "aws_ecr_repository" "this" {
  name = var.name

  image_scanning_configuration {
    scan_on_push = var.ecr_image_scan_on_push
  }

  image_tag_mutability = var.ecr_image_tag_mutability

  tags = var.tags
}
