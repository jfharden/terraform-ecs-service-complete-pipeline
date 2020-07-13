locals {
  ci_name      = "ci-${var.name}"
  ci_log_group = "ci-build-logs"
}

resource "aws_codebuild_project" "ci" {
  name          = local.ci_name
  description   = "CI build for ${var.github_url}"
  build_timeout = var.ci_build_timeout
  service_role  = aws_iam_role.this.arn

  source {
    type            = "GITHUB"
    location        = var.github_url
    git_clone_depth = var.github_clone_depth
    buildspec       = var.ci_buildspec
  }

  cache {
    type = "NO_CACHE"
  }

  source_version = "master"

  artifacts {
    type = "NO_ARTIFACTS"
  }

  logs_config {
    cloudwatch_logs {
      status      = "ENABLED"
      group_name  = local.ci_log_group
      stream_name = local.ci_name
    }
  }

  environment {
    compute_type    = var.ci_compute_type
    image           = var.ci_codebuild_image
    type            = "LINUX_CONTAINER"
    privileged_mode = var.ci_privileged_mode

    environment_variable {
      name  = "AWS_REGION"
      value = var.aws_region
    }
  }

  tags = var.tags
}

resource "aws_codebuild_webhook" "ci" {
  project_name = aws_codebuild_project.ci.name

  filter_group {
    filter {
      type    = "EVENT"
      pattern = "PULL_REQUEST_CREATED,PULL_REQUEST_REOPENED,PULL_REQUEST_UPDATED"
    }
  }
}

resource "aws_iam_role" "ci" {
  name               = local.ci_name
  assume_role_policy = data.aws_iam_policy_document.ci_assume_role.json

  tags = var.tags
}

data "aws_iam_policy_document" "ci_assume_role" {
  statement {
    sid = ""

    actions = [
      "sts:AssumeRole",
    ]

    principals {
      type        = "Service"
      identifiers = ["codebuild.amazonaws.com"]
    }

    effect = "Allow"
  }
}

resource "aws_iam_role_policy" "ci" {
  name   = local.ci_name
  role   = aws_iam_role.ci.id
  policy = data.aws_iam_policy_document.this.json
}

data "aws_iam_policy_document" "ci" {
  source_json = var.ci_codebuild_iam_policy_document_json

  statement {
    sid = "CreateLogGroup"

    effect = "Allow"

    actions = [
      "logs:CreateLogGroup",
    ]

    resources = ["*"]
  }

  statement {
    sid = "CreateLogStream"

    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    effect = "Allow"

    resources = [
      "arn:aws:logs:${var.aws_region}:${var.account_id}:log-group:${local.ci_log_group}:*",
    ]
  }
}
