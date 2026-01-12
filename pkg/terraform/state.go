package terraform

import (
	"encoding/json"
	"fmt"
	"os"
)

// StateParser parses Terraform state files
type StateParser struct {
	statePath string
	state     *State
}

// NewStateParser creates a new Terraform state parser
func NewStateParser(statePath string) *StateParser {
	return &StateParser{
		statePath: statePath,
	}
}

// loadState loads and parses the Terraform state file
func (p *StateParser) loadState() error {
	if p.state != nil {
		return nil
	}

	data, err := os.ReadFile(p.statePath)
	if err != nil {
		return fmt.Errorf("failed to read state file: %w", err)
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("failed to parse state file: %w", err)
	}

	p.state = &state
	return nil
}

// GetInstanceConfig retrieves configuration for a specific EC2 instance
func (p *StateParser) GetInstanceConfig(instanceID string) (map[string]any, error) {
	if err := p.loadState(); err != nil {
		return nil, err
	}

	// Find the EC2 instance resource
	for _, resource := range p.state.Resources {
		if resource.Type == "aws_instance" {
			for _, instance := range resource.Instances {
				if id, ok := instance.Attributes["id"].(string); ok && id == instanceID {
					return p.normalizeAttributes(instance.Attributes), nil
				}
			}
		}
	}
	return nil, fmt.Errorf("instance %s not found in Terraform state", instanceID)
}

// normalizeAttributes normalizes Terraform attributes to match AWS format
func (p *StateParser) normalizeAttributes(attrs map[string]any) map[string]any {
	normalized := make(map[string]any)

	// Copy all attributes
	for k, v := range attrs {
		normalized[k] = v
	}

	// Security groups - convert to string slice
	if sgIDs, ok := attrs["vpc_security_group_ids"]; ok {
		normalized["vpc_security_group_ids"] = convertToStringSlice(sgIDs)
	}

	// Tags - ensure proper format
	if tags, ok := attrs["tags"].(map[string]any); ok {
		normalized["tags"] = tags
	}

	// Convert float64 to int64 where appropriate
	for k, v := range normalized {
		if f, ok := v.(float64); ok {
			if f == float64(int64(f)) {
				normalized[k] = int64(f)
			}
		}
	}

	return normalized
}

// convertToStringSlice converts various types to []string
func convertToStringSlice(val any) []string {
	switch v := val.(type) {
	case []string:
		return v
	case []any:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
		return result
	case string:
		return []string{v}
	default:
		return []string{}
	}
}

// GetAllInstances returns all EC2 instances from the state
func (p *StateParser) GetAllInstances() ([]map[string]any, error) {
	if err := p.loadState(); err != nil {
		return nil, err
	}

	instances := make([]map[string]any, 0)

	for _, resource := range p.state.Resources {
		if resource.Type == "aws_instance" {
			for _, instance := range resource.Instances {
				instances = append(instances, p.normalizeAttributes(instance.Attributes))
			}
		}
	}

	return instances, nil
}

// GetInstanceIDs returns all EC2 instance IDs from the state
func (p *StateParser) GetInstanceIDs() ([]string, error) {
	instances, err := p.GetAllInstances()
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(instances))
	for _, instance := range instances {
		if id, ok := instance["id"].(string); ok {
			ids = append(ids, id)
		}
	}

	return ids, nil
}