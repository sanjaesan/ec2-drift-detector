# AWS EC2 Terraform Drift Detection Tool

A Go application that detects infrastructure drift between AWS EC2 instances and their Terraform configurations.

## Features

- ✅ Multi-attribute drift detection (instance_type, AMI, subnet, security groups, tags, etc.)
- ✅ Concurrent processing for multiple instances
- ✅ Mock mode for testing without AWS credentials
- ✅ Structured console and JSON output
- ✅ Support for nested attributes (tags, security groups)
- ✅ >70% test coverage

## Project Structure

```
ec2-drift-detector/
├── cmd/drift-detector/      # Application entry point
├── pkg/                      # Public packages
│   ├── detector/            # Drift detection logic
│   ├── aws/                 # AWS EC2 integration
│   ├── terraform/           # Terraform state parsing
│   └── reporter/            # Output formatting
├── internal/appconfig/      # Internal configuration
└── testdata/                # Test fixtures
```

## Prerequisites

- Go 1.19 or higher
- AWS credentials (for live mode)
- Terraform state file

## Installation

```bash
# Clone the repository
git clone https://github.com/sanjaesan/ec2-drift-detector.git
cd ec2-drift-detector

# Install dependencies
make install

# Build
make build
```

## Usage

### Quick Start (Mock Mode)

```bash
make run-mock
```

### With Real AWS

```bash
# Set AWS credentials
export AWS_ACCESS_KEY_ID="your-key"
export AWS_SECRET_ACCESS_KEY="your-secret"
export AWS_REGION="us-east-1"

# Run
./drift-detector \
  --instances=i-0abc123def456789 \
  --terraform-state=testdata/terraform.tfstate
```

### Concurrent Mode

```bash
./drift-detector \
  --instances=i-xxx,i-yyy,i-zzz \
  --terraform-state=testdata/terraform.tfstate \
  --concurrent
```

### JSON Output

```bash
./drift-detector \
  --instances=i-xxx \
  --terraform-state=testdata/terraform.tfstate \
  --format=json
```

## CLI Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--instances` | Comma-separated EC2 instance IDs | Required |
| `--terraform-state` | Path to Terraform state file | `terraform.tfstate` |
| `--attributes` | Attributes to check | `instance_type,ami,subnet_id,vpc_security_group_ids,tags` |
| `--mock` | Use mock data | `false` |
| `--concurrent` | Enable concurrent processing | `false` |
| `--format` | Output format (console/json) | `console` |

## Development

### Run Tests

```bash
make test
```

### Generate Coverage Report

```bash
make coverage
```

### Format Code

```bash
make fmt
```

### Run All Checks

```bash
make all
```

## Example Output

```
================================================================================
EC2 TERRAFORM DRIFT DETECTION REPORT
================================================================================

Instance: i-1234567890abcdef0
--------------------------------------------------------------------------------
Drift Detected: YES (1 attribute(s))

  1. Attribute: instance_type
     AWS Value:       "t3.medium"
     Terraform Value: "t3.small"

================================================================================
SUMMARY
================================================================================
Total Instances Checked: 1
Instances with Drift:    1
Instances with Errors:   0
Instances in Sync:       0
================================================================================
```

## Design Decisions

### Interface-Based Architecture
All major components use interfaces for testability and flexibility.

### Concurrent Processing
Uses worker pool pattern with semaphore limiting (10 concurrent requests) to respect AWS API rate limits.

### Value Comparison
Implements type-safe comparison for strings, slices, maps, and nested structures.

## Testing

```bash
# Run all tests
make test

# Run with race detection
make test-race

# View coverage
make coverage
```

Current test coverage: ~75%

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Submit a pull request

## License

MIT License

## Author

[Jubril Sanusi](mailto:sanusijubril.sj@gmail.com)

## Future Improvements

- **Google Cloud Platform (GCP)**: Support for Compute Engine instances
- **Azure**: Support for Virtual Machines
- **Multi-cloud detection**: Unified interface for all providers

## Acknowledgments

- AWS SDK for Go v2
- HashiCorp Terraform
- Go community