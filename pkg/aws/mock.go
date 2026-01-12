package aws

import (
	"context"
	"fmt"
)

type MockEC2Client struct {
	instances map[string]map[string]interface{}
}

func NewMockEC2Client() *MockEC2Client {
	return &MockEC2Client{
		instances: map[string]map[string]interface{}{
			"i-1234567890abcdef0": {
				"instance_type": "t3.medium",
				"ami":           "ami-0c55b159cbfafe1f0",
				"subnet_id":     "subnet-12345678",
				"vpc_id":        "vpc-12345678",
				"key_name":      "my-key-pair",
				"vpc_security_group_ids": []string{
					"sg-12345678",
					"sg-87654321",
				},
				"tags": map[string]interface{}{
					"Name":        "web-server-1",
					"Environment": "production",
					"ManagedBy":   "terraform",
				},
				"monitoring": "disabled",
			},
			"i-0987654321fedcba0": {
				"instance_type": "t3.large",
				"ami":           "ami-0c55b159cbfafe1f0",
				"subnet_id":     "subnet-87654321",
				"vpc_id":        "vpc-12345678",
				"key_name":      "my-key-pair",
				"vpc_security_group_ids": []string{
					"sg-12345678",
				},
				"tags": map[string]interface{}{
					"Name":        "web-server-2",
					"Environment": "staging",
					"ManagedBy":   "terraform",
				},
				"monitoring": "enabled",
			},
		},
	}
}

func (m *MockEC2Client) GetInstance(ctx context.Context, instanceID string) (map[string]interface{}, error) {
	instance, exists := m.instances[instanceID]
	if !exists {
		return nil, fmt.Errorf("instance %s not found in mock data", instanceID)
	}

	// Return a copy to prevent modifications
	return copyMap(instance), nil
}

func copyMap(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{})
	for k, v := range src {
		switch val := v.(type) {
		case map[string]interface{}:
			dst[k] = copyMap(val)
		case []string:
			copySlice := make([]string, len(val))
			copy(copySlice, val)
			dst[k] = copySlice
		default:
			dst[k] = v
		}
	}
	return dst
}