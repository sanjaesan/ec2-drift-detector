package detector

type Result struct {
	InstanceID string
	HasDrift   bool
	Drifts     []AttributeDrift
	Error      error
}

type AttributeDrift struct {
	Attribute      string
	AWSValue       any
	TerraformValue any
	Path           string // For nested attributes
}