package detector

import (
	"context"
	"testing"
)

// Mock EC2 Client for testing
type mockEC2Client struct {
	instances map[string]map[string]any
}

func (m *mockEC2Client) GetInstance(ctx context.Context, instanceID string) (map[string]any, error) {
	config, exists := m.instances[instanceID]
	if !exists {
		return nil, nil
	}
	return config, nil
}

// Mock Terraform Parser for testing
type mockTerraformParser struct {
	instances map[string]map[string]any
}

func (m *mockTerraformParser) GetInstanceConfig(instanceID string) (map[string]any, error) {
	config, exists := m.instances[instanceID]
	if !exists {
		return nil, nil
	}
	return config, nil
}

func (m *mockTerraformParser) GetAllInstances() ([]map[string]any, error) {
	instances := make([]map[string]any, 0, len(m.instances))
	for _, inst := range m.instances {
		instances = append(instances, inst)
	}
	return instances, nil
}

func (m *mockTerraformParser) GetInstanceIDs() ([]string, error) {
	ids := make([]string, 0, len(m.instances))
	for id := range m.instances {
		ids = append(ids, id)
	}
	return ids, nil
}

// Tests
func TestDetector_Detect_NoDrift(t *testing.T) {
	ec2Client := &mockEC2Client{
		instances: map[string]map[string]any{
			"i-test": {
				"instance_type": "t3.medium",
				"ami":           "ami-12345",
			},
		},
	}

	tfParser := &mockTerraformParser{
		instances: map[string]map[string]any{
			"i-test": {
				"instance_type": "t3.medium",
				"ami":           "ami-12345",
			},
		},
	}

	detector := New(ec2Client, tfParser, []string{"instance_type", "ami"})
	results, err := detector.Detect(context.Background(), []string{"i-test"})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].HasDrift {
		t.Errorf("Expected no drift")
	}
}

func TestDetector_Detect_WithDrift(t *testing.T) {
	ec2Client := &mockEC2Client{
		instances: map[string]map[string]any{
			"i-test": {
				"instance_type": "t3.large",
			},
		},
	}

	tfParser := &mockTerraformParser{
		instances: map[string]map[string]any{
			"i-test": {
				"instance_type": "t3.medium",
			},
		},
	}

	detector := New(ec2Client, tfParser, []string{"instance_type"})
	results, err := detector.Detect(context.Background(), []string{"i-test"})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !results[0].HasDrift {
		t.Errorf("Expected drift to be detected")
	}

	if len(results[0].Drifts) != 1 {
		t.Errorf("Expected 1 drift, got %d", len(results[0].Drifts))
	}
}

func TestDetector_DetectConcurrent(t *testing.T) {
	ec2Client := &mockEC2Client{
		instances: map[string]map[string]any{
			"i-test1": {"instance_type": "t3.medium"},
			"i-test2": {"instance_type": "t3.large"},
			"i-test3": {"instance_type": "t3.small"},
		},
	}

	tfParser := &mockTerraformParser{
		instances: map[string]map[string]any{
			"i-test1": {"instance_type": "t3.medium"},
			"i-test2": {"instance_type": "t3.medium"},
			"i-test3": {"instance_type": "t3.small"},
		},
	}

	detector := New(ec2Client, tfParser, []string{"instance_type"})
	results, err := detector.DetectConcurrent(context.Background(),
		[]string{"i-test1", "i-test2", "i-test3"})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}
}

func TestDetector_SecurityGroupsDrift(t *testing.T) {
	ec2Client := &mockEC2Client{
		instances: map[string]map[string]any{
			"i-test": {
				"vpc_security_group_ids": []string{"sg-123", "sg-456"},
			},
		},
	}

	tfParser := &mockTerraformParser{
		instances: map[string]map[string]any{
			"i-test": {
				"vpc_security_group_ids": []string{"sg-123"},
			},
		},
	}

	detector := New(ec2Client, tfParser, []string{"vpc_security_group_ids"})
	results, err := detector.Detect(context.Background(), []string{"i-test"})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !results[0].HasDrift {
		t.Errorf("Expected drift to be detected")
	}
}

func TestDetector_TagsDrift(t *testing.T) {
	ec2Client := &mockEC2Client{
		instances: map[string]map[string]any{
			"i-test": {
				"tags": map[string]any{
					"Environment": "production",
				},
			},
		},
	}

	tfParser := &mockTerraformParser{
		instances: map[string]map[string]any{
			"i-test": {
				"tags": map[string]any{
					"Environment": "staging",
				},
			},
		},
	}

	detector := New(ec2Client, tfParser, []string{"tags"})
	results, err := detector.Detect(context.Background(), []string{"i-test"})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !results[0].HasDrift {
		t.Errorf("Expected drift in tags")
	}
}