package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

type TestData struct {
	AwsRegion        string
	RandomName       string
	TerraformOptions *terraform.Options
	WorkingDir       string
}

func LoadTestData(t *testing.T, workingDir string) *TestData {
	return &TestData{
		AwsRegion:        test_structure.LoadString(t, workingDir, "awsRegion"),
		RandomName:       test_structure.LoadString(t, workingDir, "randomName"),
		TerraformOptions: test_structure.LoadTerraformOptions(t, workingDir),
		WorkingDir:       workingDir,
	}
}

func (testData *TestData) Save(t *testing.T) {
	test_structure.SaveString(t, testData.WorkingDir, "randomName", testData.RandomName)
	test_structure.SaveString(t, testData.WorkingDir, "awsRegion", testData.AwsRegion)
	test_structure.SaveTerraformOptions(t, testData.WorkingDir, testData.TerraformOptions)
}

func Setup(t *testing.T, workingDir string) *TestData {
	randomName := fmt.Sprintf("ecs-test-%s", strings.ToLower(random.UniqueId()))
	awsRegion := "eu-west-1"

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/terraform-ecs-service-complete-pipeline/",

		Vars: map[string]interface{}{
			"name": randomName,
			"tags": map[string]string{
				"TerraformTest": randomName,
				"Stage":         "test",
			},
		},

		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
	}

	testData := &TestData{
		AwsRegion:        awsRegion,
		RandomName:       randomName,
		TerraformOptions: terraformOptions,
		WorkingDir:       workingDir,
	}

	testData.Save(t)

	return testData
}

func TestServicePipeline(t *testing.T) {
	workingDir := filepath.Join(".terratest-working-dir", t.Name())

	defer test_structure.RunTestStage(t, "destroy", func() {
		testData := LoadTestData(t, workingDir)
		terraform.Destroy(t, testData.TerraformOptions)
	})

	test_structure.RunTestStage(t, "init_apply", func() {
		testData := Setup(t, workingDir)

		terraform.InitAndApply(t, testData.TerraformOptions)
	})

	validators := map[string]func(*testing.T, *TestData){
		"ValidateECRRepository": ValidateECR,
	}

	test_structure.RunTestStage(t, "validate", func() {
		testData := LoadTestData(t, workingDir)

		for name, validator := range validators {
			t.Run(name, func(t *testing.T) {
				validator(t, testData)
			})
		}
	})
}

func ValidateECR(t *testing.T, testData *TestData) {
	// Check the ARN is as expected
	accountID := aws.GetAccountId(t)
	expectedEcrArn := fmt.Sprintf("arn:aws:ecr:eu-west-1:%s:repository/%s", accountID, testData.RandomName)
	ecrArn := terraform.Output(t, testData.TerraformOptions, "ecr_repository_arn")

	assert.Equal(t, expectedEcrArn, ecrArn, "The ARN of the ecr repository is incorrect")

	// Check the repo has image scanning enabled
	imageScanningEnabled := terraform.Output(t, testData.TerraformOptions, "ecr_image_scan_on_push")
	assert.Equal(t, "true", imageScanningEnabled, "Image Scanning is not enabled")

	// Check image tag mutability
	imageTagMutability := terraform.Output(t, testData.TerraformOptions, "ecr_image_tag_mutability")
	assert.Equal(t, "IMMUTABLE", imageTagMutability, "Image tags are mutable")

	// TODO: Need to add test for tags
}
