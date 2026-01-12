package reporter

import (
	"encoding/json"
	"fmt"

	"github.com/sanjaesan/ec2-drift-detector/pkg/detector"
)

type JSONReporter struct{}

func NewJSONReporter() *JSONReporter {
	return &JSONReporter{}
}

func (r *JSONReporter) Report(results []detector.Result) {
	output := struct {
		Results []detector.Result `json:"results"`
		Summary struct {
			Total      int `json:"total"`
			WithDrift  int `json:"with_drift"`
			WithErrors int `json:"with_errors"`
			InSync     int `json:"in_sync"`
		} `json:"summary"`
	}{
		Results: results,
	}

	for _, result := range results {
		output.Summary.Total++
		if result.Error != nil {
			output.Summary.WithErrors++
		} else if result.HasDrift {
			output.Summary.WithDrift++
		} else {
			output.Summary.InSync++
		}
	}

	jsonBytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Printf("Error formatting JSON: %v\n", err)
		return
	}

	fmt.Println(string(jsonBytes))
}