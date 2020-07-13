package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	awsSdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/aws/aws-sdk-go/service/ecr"

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
		"ValidateCodebuildCI":   ValidateCodebuildCI,
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

func ValidateCodebuildCI(t *testing.T, testData *TestData) {
	// Start build
	&codebuild.startBuildInput{
		ProjectName: fmt.Sprintf("ci-%s", randomName),
		// SourceVersion: "master",
	}
	// Wait for completion
	// Check it was successful
}

func ValidateCodebuildCiWebhook(t *testing.T, testData *TestData) {
	// Make sure there is a webhook
	// Make sure the webhook is sensible
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

	// Check the tags are correct
	awsSession, err := session.NewSession(
		&awsSdk.Config{
			Region: awsSdk.String(testData.AwsRegion),
		},
	)

	actualTags, err := getECRTags(ecrArn, awsSession)
	if err != nil {
		t.Errorf("Couldn't request tags for ECR repository: %s", err)
	}

	expectedTags := map[string]string{
		"TerraformTest": testData.RandomName,
		"Stage":         "test",
	}

	assertEcrTagsMatch(t, expectedTags, actualTags)
}

func getECRTags(ecrArn string, awsSession *session.Session) ([]*ecr.Tag, error) {
	ecrService := ecr.New(awsSession)

	tagRequest, tagResponse := ecrService.ListTagsForResourceRequest(
		&ecr.ListTagsForResourceInput{
			ResourceArn: awsSdk.String(ecrArn),
		},
	)

	err := tagRequest.Send()
	if err != nil {
		return nil, err
	}

	return tagResponse.Tags, nil
}

func assertEcrTagsMatch(t *testing.T, expected map[string]string, actual []*ecr.Tag) {
	equalNumberOfTags := assert.Equal(t,
		len(expected),
		len(actual),
		fmt.Sprintf(
			"Number of tags expected [%v] does match actual number of tags [%v]. Expected: %v Actual: %v",
			len(expected),
			len(actual),
			expected,
			actual,
		),
	)

	if !equalNumberOfTags {
		return
	}

	for _, tag := range actual {
		actualKey := awsSdk.StringValue(tag.Key)
		actualValue := awsSdk.StringValue(tag.Value)

		expectedValue, tagInExpected := expected[actualKey]
		if !tagInExpected {
			t.Errorf("Error: Tag key %s in actual, but not expected\n", actualKey)
			return
		}

		if !(expectedValue == actualValue) {
			t.Errorf("Error: Tag with key %s has value %s but we expected %s", actualKey, actualValue, expectedValue)
			return
		}
	}
}
