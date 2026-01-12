package terraform

// State represents the structure of a Terraform state file
type State struct {
	Version          int            `json:"version"`
	TerraformVersion string         `json:"terraform_version"`
	Resources        []Resource     `json:"resources"`
	Outputs          map[string]any `json:"outputs,omitempty"`
}

// Resource represents a resource in Terraform state
type Resource struct {
	Mode      string             `json:"mode"`
	Type      string             `json:"type"`
	Name      string             `json:"name"`
	Provider  string             `json:"provider"`
	Instances []ResourceInstance `json:"instances"`
}

// ResourceInstance represents an instance of a resource
type ResourceInstance struct {
	SchemaVersion int            `json:"schema_version"`
	Attributes    map[string]any `json:"attributes"`
	Private       string         `json:"private,omitempty"`
	Dependencies  []string       `json:"dependencies,omitempty"`
}
