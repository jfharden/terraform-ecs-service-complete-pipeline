variable "name" {
  type        = string
  description = "Name to give to all resources"
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
