# LoveCMS Testing Documentation

This document describes the comprehensive test suite for LoveCMS, including unit tests, integration tests, benchmarks, and test helpers.

## Test Structure

The test suite is organized into several categories:

### Unit Tests
- **Application Layer** (`internal/application/*_test.go`)
  - `models_test.go` - Tests for Post and FrontMatter data structures
  - `management_test.go` - Tests for post management functions (CRUD operations)
  - `directory_picker_test.go` - Tests for the interactive directory picker
  - `benchmark_test.go` - Performance benchmarks for core functions

- **Handlers Layer** (`internal/handlers/*_test.go`)
  - `home_test.go` - Tests for home page handler
  - `post_test.go` - Tests for post-related handlers (view, create, edit)

- **Main Package** (`cmd/*_test.go`)
  - `main_test.go` - Tests for configuration and utility functions
  - `integration_test.go` - End-to-end integration tests

### Test Helpers
- `internal/testhelpers/helpers.go` - Common utilities for testing across packages

## Running Tests

### Quick Test Run
```bash
make test
```

### Test with Coverage Report
```bash
make test-coverage
```
This generates an HTML coverage report at `coverage.html`.

### Short Tests Only
```bash
make test-short
```
Skips integration tests and long-running tests.

### Benchmark Tests
```bash
make benchmark
```

### Manual Test Commands
```bash
# Run all tests with race detection
go test -race -v -timeout 30s ./...

# Run tests for specific package
go test -v ./internal/application/
go test -v ./internal/handlers/
go test -v ./cmd/

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
go test -bench=. -benchmem ./...
```

## Test Coverage

Current test coverage by package:
- **Application Package**: 81.3% coverage
- **Handlers Package**: 86.8% coverage  
- **CMD Package**: 6.5% coverage (mostly main function)

## Test Categories

### 1. Unit Tests

#### Application Layer Tests
- **Post Management**: Creating, reading, updating posts
- **File Operations**: File system interactions
- **Slug Generation**: URL-friendly slug creation
- **Directory Picker**: Interactive directory selection
- **Error Handling**: Invalid input and edge cases

#### Handler Tests
- **HTTP Responses**: Status codes, redirects, error handling
- **Form Processing**: POST data validation and processing
- **Configuration**: Site path validation and setup flow
- **Template Rendering**: Response content validation

#### Configuration Tests
- **Cross-platform Paths**: Config directory detection
- **Environment Handling**: Different OS behaviors

### 2. Integration Tests

Integration tests verify end-to-end functionality:
- Complete post lifecycle (create, read, update)
- HTTP server integration
- File system operations
- Configuration persistence

### 3. Benchmark Tests

Performance tests for critical operations:
- `BenchmarkGetPosts` - Loading multiple posts
- `BenchmarkUpdatePost` - Saving post content
- `BenchmarkCreatePost` - Creating new posts

Current benchmark results (Apple M4):
- GetPosts: ~208µs per operation
- UpdatePost: ~31µs per operation  
- CreatePost: ~86µs per operation

### 4. Test Helpers

The `testhelpers` package provides utilities for:
- Creating temporary test sites
- Generating sample posts
- Asserting post equality
- File existence validation
- Content verification

## Test Data Management

Tests use temporary directories and files that are automatically cleaned up:
- All test files are created in `os.TempDir()`
- Cleanup functions ensure no test artifacts remain
- Tests are isolated and can run in parallel

## Writing New Tests

### Test Naming Conventions
- Test functions: `TestFunctionName`
- Benchmark functions: `BenchmarkFunctionName`
- Helper functions: `helperFunctionName` (with `t.Helper()`)

### Best Practices
1. Use table-driven tests for multiple scenarios
2. Include both positive and negative test cases
3. Test error conditions and edge cases
4. Use temporary directories for file operations
5. Clean up resources with defer statements
6. Use `t.Helper()` in helper functions

### Example Test Structure
```go
func TestExampleFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   "test",
            want:    "expected",
            wantErr: false,
        },
        // Add more test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ExampleFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ExampleFunction() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("ExampleFunction() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Continuous Integration

The test suite is designed to run in CI environments:
- No external dependencies required
- All tests use temporary directories
- Race detection enabled
- Timeout protection (30s default)
- Cross-platform compatibility

## Troubleshooting

### Common Issues
1. **File Permission Errors**: Ensure test runner has write access to temp directories
2. **Race Conditions**: Use `-race` flag to detect data races
3. **Test Timeouts**: Increase timeout with `-timeout` flag if needed
4. **Coverage Issues**: Exclude generated files from coverage reports

### Debug Tests
```bash
# Run with verbose output
go test -v ./...

# Run specific test
go test -run TestSpecificFunction ./...

# Debug with print statements (avoid in committed code)
t.Logf("Debug: value = %v", someValue)
```

## Future Enhancements

Planned test improvements:
- [ ] More comprehensive integration tests
- [ ] Property-based testing for complex operations
- [ ] Performance regression tests
- [ ] End-to-end HTTP tests with actual server
- [ ] Template rendering validation tests
- [ ] Database/persistence layer tests (if added)
