package reporter

import "github.com/sanjaesan/ec2-drift-detector/pkg/detector"

// Reporter interface for reporting drift results
type Reporter interface {
	Report(results []detector.Result)
}