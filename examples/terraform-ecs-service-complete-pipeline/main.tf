terraform {
  required_version = ">= 0.12"
}

module "service_pipeline" {
  source = "../../"

  name = var.name
  tags = var.tags

  github_url = "https://github.com/jfharden/template-php-docker-site.git"
}
