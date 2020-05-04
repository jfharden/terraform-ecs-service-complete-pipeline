terraform {
  required_version = ">= 0.12"
}

module "service_pipeline" {
  source = "../../"

  name = var.name
  tags = var.tags
}
