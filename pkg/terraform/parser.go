package terraform

// Parser interface for parsing Terraform configurations
type Parser interface {
	GetInstanceConfig(instanceID string) (map[string]any, error)
	GetAllInstances() ([]map[string]any, error)
	GetInstanceIDs() ([]string, error)
}