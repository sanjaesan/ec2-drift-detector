package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type AWSEC2Client struct {
	client *ec2.Client
}

func NewAWSEC2Client(client *ec2.Client) *AWSEC2Client {
	return &AWSEC2Client{client: client}
}

func (c *AWSEC2Client) GetInstance(ctx context.Context, instanceID string) (map[string]interface{}, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}

	result, err := c.client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instance: %w", err)
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("instance %s not found", instanceID)
	}

	instance := result.Reservations[0].Instances[0]
	return instanceToMap(instance), nil
}

func instanceToMap(instance types.Instance) map[string]interface{} {
	config := make(map[string]interface{})

	// Basic attributes
	config["instance_type"] = string(instance.InstanceType)

	if instance.ImageId != nil {
		config["ami"] = *instance.ImageId
	}

	if instance.SubnetId != nil {
		config["subnet_id"] = *instance.SubnetId
	}

	if instance.VpcId != nil {
		config["vpc_id"] = *instance.VpcId
	}

	if instance.KeyName != nil {
		config["key_name"] = *instance.KeyName
	}

	if instance.PrivateIpAddress != nil {
		config["private_ip"] = *instance.PrivateIpAddress
	}

	if instance.PublicIpAddress != nil {
		config["public_ip"] = *instance.PublicIpAddress
	}
	// Security groups
	securityGroups := make([]string, 0, len(instance.SecurityGroups))
	for _, sg := range instance.SecurityGroups {
		if sg.GroupId != nil {
			securityGroups = append(securityGroups, *sg.GroupId)
		}
	}
	if len(securityGroups) > 0 {
		config["vpc_security_group_ids"] = securityGroups
	}

	// Tags
	if len(instance.Tags) > 0 {
		tags := make(map[string]interface{})
		for _, tag := range instance.Tags {
			if tag.Key != nil && tag.Value != nil {
				tags[*tag.Key] = *tag.Value
			}
		}
		config["tags"] = tags
	}

	// Monitoring
	if instance.Monitoring != nil {
		config["monitoring"] = string(instance.Monitoring.State)
	}

	// IAM instance profile
	if instance.IamInstanceProfile != nil && instance.IamInstanceProfile.Arn != nil {
		config["iam_instance_profile"] = *instance.IamInstanceProfile.Arn
	}

	return config
}