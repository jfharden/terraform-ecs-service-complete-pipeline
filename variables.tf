variable "name" {
  type        = string
  description = "Name to apply to all the resources."
}

variable "ecr_image_scan_on_push" {
  type = bool
  description = "Enable ECR image scanning on push"
  default = true
}

variable "ecr_image_tag_mutability" {
  type = string
  description = "The tag mutability setting for the repository. Must be one of: MUTABLE or IMMUTABLE"
  default = "IMMUTABLE"
}

variable "tags" {
  type        = map
  description = "Tags to apply to all resources which can be tagged"
  default     = {}
}
