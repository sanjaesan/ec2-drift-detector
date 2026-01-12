package appconfig

// Config holds application configuration
type Config struct {
	TerraformStateFile string
	InstanceIDs        []string
	Attributes         []string
	UseMockData        bool
	Concurrent         bool
	OutputFormat       string
}