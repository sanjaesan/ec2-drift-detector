# Architecture Documentation

## Overview

The EC2 Drift Detector is designed using clean architecture principles with clear separation of concerns. The system follows a layered architecture that promotes testability, maintainability, and scalability.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                       CLI Layer                              │
│                  (cmd/drift-detector)                        │
│              • Flag Parsing                                  │
│              • Dependency Injection                          │
│              • Orchestration                                 │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                         │
│                    (pkg/detector)                            │
│              • Drift Detection Logic                         │
│              • Comparison Algorithms                         │
│              • Concurrency Management                        │
└─────┬──────────────┴──────────────┬─────────────────────────┘
      │                              │
      ▼                              ▼
┌──────────────────┐        ┌──────────────────┐
│  External APIs   │        │  Data Sources    │
│   (pkg/aws)      │        │ (pkg/terraform)  │
│                  │        │                  │
│ • EC2 Client     │        │ • State Parser   │
│ • Data Mapping   │        │ • Normalization  │
└──────────────────┘        └──────────────────┘
      │                              │
      └──────────────┬───────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                        │
│                    (pkg/reporter)                            │
│              • Console Output                                │
│              • JSON Output                                   │
│              • Result Formatting                             │
└─────────────────────────────────────────────────────────────┘
```

## Component Architecture

### 1. CLI Layer (`cmd/drift-detector`)

**Responsibility**: Application entry point and configuration

**Components**:
- `main()`: Entry point
- `parseFlags()`: Command-line argument parsing
- `hasDrift()`: Result aggregation

**Dependencies**:
- All `pkg/*` packages
- `internal/appconfig`
- AWS SDK v2 (for configuration)

**Design Pattern**: Dependency Injection

```go
// Dependency injection at startup
ec2Client := aws.NewAWSEC2Client(...)
tfParser := terraform.NewStateParser(...)
detector := detector.New(ec2Client, tfParser, attributes)
```

### 2. Detection Layer (`pkg/detector`)

**Responsibility**: Core business logic for drift detection

**Components**:

#### detector.go
- `Detector`: Main detection engine
- `Detect()`: Sequential processing
- `DetectConcurrent()`: Parallel processing
- `detectSingleInstance()`: Single instance analysis
- `compareAttribute()`: Attribute comparison

#### compare.go
- `valuesEqual()`: Type-safe value comparison
- `getNestedValue()`: Nested attribute access
- `splitPath()`: Path parsing
- Helper comparison functions

#### types.go
- `Result`: Detection result structure
- `AttributeDrift`: Drift information

**Design Patterns**:
- Strategy Pattern (for comparison)
- Worker Pool Pattern (for concurrency)
- Builder Pattern (for detector construction)

**Concurrency Model**:
```go
// Semaphore-based rate limiting
semaphore := make(chan struct{}, 10) // Max 10 concurrent
for each instance {
    go func() {
        semaphore <- struct{}{}        // Acquire
        defer func() { <-semaphore }() // Release
        detectInstance()
    }()
}
```

### 3. AWS Integration Layer (`pkg/aws`)

**Responsibility**: AWS EC2 API interaction and data transformation

**Components**:

#### client.go
- `EC2Client`: Interface defining contract

#### ec2.go
- `AWSEC2Client`: Real AWS implementation
- `GetInstance()`: Fetch instance data
- `instanceToMap()`: Transform AWS types to comparable format

#### mock.go
- `MockEC2Client`: Test implementation with sample data
- `GetInstance()`: Return mock data
- `copyMap()`: Deep copy for isolation

**Design Patterns**:
- Interface Segregation (single method interface)
- Adapter Pattern (AWS SDK → internal format)
- Mock Object Pattern (for testing)

**Data Flow**:
```
AWS API Response (types.Instance)
         ↓
instanceToMap() [Normalization]
         ↓
map[string]interface{} [Generic format]
         ↓
Detector [Comparison]
```

### 4. Terraform Integration Layer (`pkg/terraform`)

**Responsibility**: Terraform state file parsing and normalization

**Components**:

#### parser.go
- `Parser`: Interface for state access

#### state.go
- `StateParser`: JSON state file parser
- `GetInstanceConfig()`: Extract instance configuration
- `normalizeAttributes()`: Normalize to AWS format
- `convertToStringSlice()`: Type conversion helper

#### types.go
- `State`: Top-level state structure
- `Resource`: Resource representation
- `ResourceInstance`: Instance data

**Design Patterns**:
- Facade Pattern (simplifies state access)
- Adapter Pattern (Terraform → internal format)
- Lazy Loading (state loaded on first access)

**Normalization Strategy**:
```
Terraform State (JSON)
         ↓
Parse to structs
         ↓
normalizeAttributes()
    • Security groups → []string
    • Tags → map[string]interface{}
    • float64 → int64 (where appropriate)
         ↓
map[string]interface{} [Generic format]
         ↓
Detector [Comparison]
```

### 5. Reporting Layer (`pkg/reporter`)

**Responsibility**: Output formatting and presentation

**Components**:

#### reporter.go
- `Reporter`: Interface for output strategies

#### console.go
- `ConsoleReporter`: Human-readable terminal output
- `Report()`: Format and print results
- `formatValue()`: Pretty-print values

#### json.go
- `JSONReporter`: Machine-readable JSON output
- `Report()`: Serialize results

**Design Patterns**:
- Strategy Pattern (multiple output formats)
- Template Method (common reporting flow)

**Output Flow**:
```
[]detector.Result
      ↓
Reporter.Report()
      ↓
   Format
      ↓
    Output
```

### 6. Configuration Layer (`internal/appconfig`)

**Responsibility**: Application configuration structure

**Components**:
- `Config`: Centralized configuration struct

**Design Pattern**: Data Transfer Object (DTO)

## Data Flow

### Complete Request Flow

```
1. CLI Input
   ↓
2. Parse Flags → Config
   ↓
3. Initialize Dependencies
   • EC2Client (real or mock)
   • TerraformParser
   • Detector
   ↓
4. Detect() or DetectConcurrent()
   ↓
5. For each instance:
   a. Fetch AWS config (EC2Client.GetInstance)
   b. Fetch TF config (Parser.GetInstanceConfig)
   c. Compare attributes
   d. Build Result
   ↓
6. Aggregate Results
   ↓
7. Report (Console or JSON)
   ↓
8. Exit (code 0 or 1)
```

### Data Transformation Pipeline

```
AWS EC2 Instance          Terraform State
(types.Instance)          (JSON)
        ↓                        ↓
   instanceToMap()          Parse + Normalize
        ↓                        ↓
map[string]interface{}    map[string]interface{}
        ↓                        ↓
        └────────┬───────────────┘
                 ↓
          valuesEqual()
                 ↓
         AttributeDrift (if different)
                 ↓
             Result
                 ↓
            Reporter
                 ↓
        Formatted Output
```

## Design Principles

### 1. Interface-Based Design

All external dependencies are defined as interfaces:

```go
// Enables testing and flexibility
type EC2Client interface {
    GetInstance(ctx, id) (map[string]interface{}, error)
}

type Parser interface {
    GetInstanceConfig(id) (map[string]interface{}, error)
}

type Reporter interface {
    Report(results []detector.Result)
}
```

**Benefits**:
- Easy to mock for testing
- Can swap implementations
- Loose coupling between components

### 2. Dependency Injection

Dependencies are injected rather than created internally:

```go
// Constructor injection
func New(ec2Client aws.EC2Client, 
         tfParser terraform.Parser, 
         attributes []string) *Detector {
    return &Detector{
        ec2Client:  ec2Client,
        tfParser:   tfParser,
        attributes: attributes,
    }
}
```

**Benefits**:
- Testable without external dependencies
- Configurable behavior
- Clear dependencies

### 3. Single Responsibility Principle

Each package has one clear purpose:
- `detector`: Drift detection logic only
- `aws`: AWS interaction only
- `terraform`: Terraform parsing only
- `reporter`: Output formatting only

### 4. Open/Closed Principle

System is open for extension, closed for modification:

```go
// Add new reporter without changing existing code
type HTMLReporter struct{}
func (r *HTMLReporter) Report(results []detector.Result) {
    // HTML output implementation
}
```

### 5. Dependency Inversion Principle

High-level modules don't depend on low-level modules:

```
detector (high-level)
    ↓ depends on
EC2Client interface (abstraction)
    ↑ implemented by
AWSEC2Client (low-level)
```

## Concurrency Architecture

### Worker Pool Pattern

```go
// Semaphore limits concurrent workers
semaphore := make(chan struct{}, 10)

// WaitGroup coordinates completion
var wg sync.WaitGroup

// Pre-allocated results (avoid race conditions)
results := make([]Result, len(instanceIDs))

for i, id := range instanceIDs {
    wg.Add(1)
    go func(idx int, instanceID string) {
        defer wg.Done()
        
        // Acquire semaphore
        semaphore <- struct{}{}
        defer func() { <-semaphore }()
        
        // Safe: each goroutine writes to unique index
        results[idx] = detectInstance(instanceID)
    }(i, id)
}

wg.Wait()
```

**Why this design?**
- **Semaphore**: Limits concurrent API calls (respects AWS rate limits)
- **Pre-allocated slice**: Avoids race conditions on append
- **WaitGroup**: Ensures all goroutines complete
- **Fixed index**: Each goroutine writes to its own slot

### Race Condition Prevention

❌ **Wrong** (race condition):
```go
results := []Result{}
for _, id := range ids {
    go func(instanceID string) {
        r := detectInstance(instanceID)
        results = append(results, r) // RACE!
    }(id)
}
```

✅ **Correct** (no race):
```go
results := make([]Result, len(ids))
for i, id := range ids {
    go func(idx int, instanceID string) {
        results[idx] = detectInstance(instanceID) // Safe
    }(i, id)
}
```

## Error Handling Strategy

### Per-Instance Error Handling

Errors are captured per-instance, not fail-fast:

```go
result := Result{InstanceID: id}

awsConfig, err := ec2Client.GetInstance(ctx, id)
if err != nil {
    result.Error = fmt.Errorf("failed to get AWS instance: %w", err)
    return result // Continue with other instances
}
```

**Benefits**:
- Partial results are useful
- User can fix issues incrementally
- Better UX than failing on first error

### Error Wrapping

Errors are wrapped with context:

```go
return fmt.Errorf("failed to get AWS instance: %w", err)
```

**Benefits**:
- Preserves original error
- Adds context
- Enables error inspection

## Testing Strategy

### Unit Test Architecture

```
Production Code          Test Code
──────────────          ──────────
detector.go     ←────── detector_test.go
    ↓ uses                    ↓ uses
EC2Client       ←────── mockEC2Client
Parser          ←────── mockTerraformParser
```

### Test Doubles

**Mock Objects** (in test files):
```go
type mockEC2Client struct {
    instances map[string]map[string]interface{}
}

func (m *mockEC2Client) GetInstance(...) {
    // Return controlled test data
}
```

**Benefits**:
- No AWS credentials needed
- Fast tests
- Controlled scenarios

### Table-Driven Tests

```go
tests := []struct{
    name     string
    input    interface{}
    expected interface{}
}{
    {"case 1", input1, expected1},
    {"case 2", input2, expected2},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // Test logic
    })
}
```

## Performance Characteristics

### Time Complexity

| Operation | Complexity | Notes |
|-----------|-----------|-------|
| Single instance detection | O(a) | a = number of attributes |
| Sequential detection | O(n * a) | n = instances, a = attributes |
| Concurrent detection | O((n/w) * a) | w = workers (10) |
| Attribute comparison | O(1) to O(m) | m = nested depth |

### Space Complexity

| Component | Complexity | Notes |
|-----------|-----------|-------|
| Results storage | O(n * d) | n = instances, d = drifts |
| Terraform state | O(r) | r = resources in state |
| Goroutines | O(min(n, 10)) | Limited by semaphore |

### Scalability

**Current limits**:
- Concurrent workers: 10 (AWS rate limit)
- Instances per run: Limited by AWS API pagination
- State file size: Limited by available memory

**Scaling options**:
- Batch processing for 1000+ instances
- Streaming state file parser for huge states
- Distributed processing for enterprise scale

## Extension Points

### Adding New Attributes

1. Update `pkg/aws/ec2.go`:
```go
func instanceToMap(instance types.Instance) map[string]interface{} {
    // Add new attribute
    config["new_attribute"] = instance.NewAttribute
}
```

2. Update `pkg/terraform/state.go` if normalization needed

3. Add test cases

### Adding New Cloud Providers

1. Create `pkg/gcp/` or `pkg/azure/`
2. Implement provider-specific client
3. Create interface compatible with `EC2Client`
4. Update `cmd/drift-detector/main.go` to support provider flag

### Adding New Output Formats

1. Create `pkg/reporter/html.go` or `pkg/reporter/slack.go`
2. Implement `Reporter` interface
3. Register in `main.go`

```go
// pkg/reporter/html.go
type HTMLReporter struct{}

func (r *HTMLReporter) Report(results []detector.Result) {
    // Generate HTML
}

// cmd/drift-detector/main.go
if cfg.OutputFormat == "html" {
    rep = reporter.NewHTMLReporter()
}
```

## Security Considerations

### AWS Credentials

- Never hardcode credentials
- Use AWS SDK credential chain
- Support IAM roles for EC2/ECS
- Environment variables for local dev

### State File Access

- State files may contain sensitive data
- Ensure proper file permissions
- Don't log state file contents
- Consider encryption at rest

### Error Messages

- Don't expose sensitive data in errors
- Sanitize error messages before display
- Log detailed errors separately

## Monitoring and Observability

### Logging

Current logging points:
- Client type (mock vs real)
- Processing mode (sequential vs concurrent)
- Instance count
- Errors (via stderr)

### Metrics to Consider

For production deployment:
- Instances checked per run
- Drifts detected
- API call duration
- Success/failure rate
- Concurrent workers utilized

### Tracing

For distributed systems:
- Add OpenTelemetry spans
- Trace request flow
- Measure component latency

## Future Architecture Improvements

### 1. Plugin System

Allow custom comparators:
```go
type AttributeComparator interface {
    Compare(aws, tf interface{}) *AttributeDrift
}

detector.RegisterComparator("tags", NewTagsComparator())
```

### 2. Result Caching

Cache results to avoid redundant API calls:
```go
type CachedDetector struct {
    detector *Detector
    cache    map[string]Result
    ttl      time.Duration
}
```

### 3. Streaming Results

Stream results as they're available:
```go
func (d *Detector) DetectStream(ctx, ids) <-chan Result {
    results := make(chan Result)
    go func() {
        // Stream results
    }()
    return results
}
```

### 4. Web API

Expose as REST API:
```
POST /api/v1/detect
{
  "instances": ["i-123", "i-456"],
  "attributes": ["instance_type", "ami"]
}
```

## Conclusion

The EC2 Drift Detector is built on solid software engineering principles:

- **Clean Architecture**: Clear separation of concerns
- **SOLID Principles**: Maintainable and extensible
- **Interface-Based**: Testable and flexible
- **Concurrent**: Efficient for multiple instances
- **Well-Tested**: High confidence in functionality

The architecture supports growth from a simple CLI tool to a comprehensive infrastructure drift detection platform.