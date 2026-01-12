package detector

import (
	"context"
	"fmt"
	"sync"

	"github.com/sanjaesan/ec2-drift-detector/pkg/aws"
	"github.com/sanjaesan/ec2-drift-detector/pkg/terraform"
)

type Detector struct {
	ec2Client  aws.EC2Client
	tfParser   terraform.Parser
	attributes []string
}

func New(ec2Client aws.EC2Client, tfParser terraform.Parser, attributes []string) *Detector {
	return &Detector{
		ec2Client:  ec2Client,
		tfParser:   tfParser,
		attributes: attributes,
	}
}

func (d *Detector) Detect(ctx context.Context, instanceIDs []string) ([]Result, error) {
	results := make([]Result, 0, len(instanceIDs))

	for _, instanceID := range instanceIDs {
		result := d.detectSingleInstance(ctx, instanceID)
		results = append(results, result)
	}

	return results, nil
}

func (d *Detector) DetectConcurrent(ctx context.Context, instanceIDs []string) ([]Result, error) {
	results := make([]Result, len(instanceIDs))
	var wg sync.WaitGroup

	// Use a semaphore to limit concurrent API calls
	semaphore := make(chan struct{}, 10)

	for i, instanceID := range instanceIDs {
		wg.Add(1)
		go func(idx int, id string) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			results[idx] = d.detectSingleInstance(ctx, id)
		}(i, instanceID)
	}

	wg.Wait()
	return results, nil
}

func (d *Detector) detectSingleInstance(ctx context.Context, instanceID string) Result {
	result := Result{
		InstanceID: instanceID,
		Drifts:     make([]AttributeDrift, 0),
	}

	// Get AWS configuration
	awsConfig, err := d.ec2Client.GetInstance(ctx, instanceID)
	if err != nil {
		result.Error = fmt.Errorf("failed to get AWS instance: %w", err)
		return result
	}

	// Get Terraform configuration
	tfConfig, err := d.tfParser.GetInstanceConfig(instanceID)
	if err != nil {
		result.Error = fmt.Errorf("failed to get Terraform config: %w", err)
		return result
	}

	// Compare attributes
	for _, attr := range d.attributes {
		drift := d.compareAttribute(attr, awsConfig, tfConfig)
		if drift != nil {
			result.Drifts = append(result.Drifts, *drift)
			result.HasDrift = true
		}
	}

	return result
}

func (d *Detector) compareAttribute(attr string, awsConfig, tfConfig map[string]interface{}) *AttributeDrift {
	awsValue := getNestedValue(awsConfig, attr)
	tfValue := getNestedValue(tfConfig, attr)

	if !valuesEqual(awsValue, tfValue) {
		return &AttributeDrift{
			Attribute:      attr,
			AWSValue:       awsValue,
			TerraformValue: tfValue,
			Path:           attr,
		}
	}

	return nil
}