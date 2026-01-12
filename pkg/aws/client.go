package aws

import (
	"context"
)

type EC2Client interface {
	GetInstance(ctx context.Context, instanceID string) (map[string]interface{}, error)
}