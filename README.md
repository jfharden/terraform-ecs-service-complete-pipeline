# terraform-ecs-service-complete-pipeline

This repo aims to provide a complete end to end pipeline which will deploy another github repository into an existing
ecs cluster and include a full ci and deploy pipeline.

## What's provided

1. ECS task with iam execution and task roles
2. ECS service to run the task
3. ECR repo to store the docker images
4. Route53 record in an existing hosted zone with ACM provided SSL certificate
5. ALB target group for the service to register with
6. Codebuild to run unit tests on the github repo, triggered on every PR open, or change of a PR
7. Codepipeline to deploy to the ECS repo
8. Notifications to an SNS topic when the service is down

## What's required

1. Existing ECS cluster with registered EC2 instances
2. Security groups allowing traffic on ephemeral ports from ALB to EC2 instances
3. Seucrity groups allowing port 443 into the ALB
4. Existing route53 hosted zone to create the subdomain under
5. Existing ALB
6. Existing VPC
7. Existing SNS topic to notify when health checks fail

## Requirements

1. go 1.13+ for running the tests
2. terraform 0.12 to use the module
3. AWS account to deploy

## Project layout

The project is broken up into the following folders

* examples - full featured examples which the tests use to deploy and destroy the module
* module - terraform module
* tests - terratest tests of the module
