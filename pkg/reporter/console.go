package reporter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sanjaesan/ec2-drift-detector/pkg/detector"
)

type ConsoleReporter struct{}

func NewConsoleReporter() *ConsoleReporter {
	return &ConsoleReporter{}
}

// Report prints drift results to console
func (r *ConsoleReporter) Report(results []detector.Result) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("EC2 TERRAFORM DRIFT DETECTION REPORT")
	fmt.Println(strings.Repeat("=", 80))

	totalDrift := 0
	totalErrors := 0

	for _, result := range results {
		fmt.Printf("\nInstance: %s\n", result.InstanceID)
		fmt.Println(strings.Repeat("-", 80))

		if result.Error != nil {
			fmt.Printf("Error: %v\n", result.Error)
			totalErrors++
			continue
		}

		if result.HasDrift {
			fmt.Printf("Drift Detected: YES (%d attribute(s))\n\n", len(result.Drifts))
			totalDrift++

			for i, drift := range result.Drifts {
				fmt.Printf("  %d. Attribute: %s\n", i+1, drift.Attribute)
				fmt.Printf("     AWS Value:       %s\n", formatValue(drift.AWSValue))
				fmt.Printf("     Terraform Value: %s\n", formatValue(drift.TerraformValue))

				if i < len(result.Drifts)-1 {
					fmt.Println()
				}
			}
		} else {
			fmt.Println("Drift Detected: NO")
			fmt.Println("All checked attributes match between AWS and Terraform")
		}
	}

	// Summary
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("SUMMARY")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Total Instances Checked: %d\n", len(results))
	fmt.Printf("Instances with Drift:    %d\n", totalDrift)
	fmt.Printf("Instances with Errors:   %d\n", totalErrors)
	fmt.Printf("Instances in Sync:       %d\n", len(results)-totalDrift-totalErrors)
	fmt.Println(strings.Repeat("=", 80) + "\n")
}

// formatValue formats a value for display
func formatValue(val any) string {
	if val == nil {
		return "<nil>"
	}

	switch v := val.(type) {
	case string:
		return fmt.Sprintf(`"%s"`, v)
	case []string:
		return fmt.Sprintf("[%s]", strings.Join(quoteStrings(v), ", "))
	case []any:
		strs := make([]string, len(v))
		for i, item := range v {
			strs[i] = formatValue(item)
		}
		return fmt.Sprintf("[%s]", strings.Join(strs, ", "))
	case map[string]any:
		jsonBytes, err := json.MarshalIndent(v, "     ", "  ")
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return "\n     " + string(jsonBytes)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// quoteStrings adds quotes around each string in a slice
func quoteStrings(strs []string) []string {
	quoted := make([]string, len(strs))
	for i, s := range strs {
		quoted[i] = fmt.Sprintf(`"%s"`, s)
	}
	return quoted
}