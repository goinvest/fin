# -*- Justfile -*-

coverage_file := "coverage.out"

# List the available justfile recipes.
@default:
  just --list --unsorted

# Format and vet Go code. Runs before tests.
[group('test')]
check:
	go fmt ./...
	go vet ./...

# Lint code using staticcheck.
[group('test')]
lint: check
	staticcheck -f stylish ./...

# Run the unit tests.
[group('test')]
unit *FLAGS: check
  go test ./... -cover -vet=off -race {{FLAGS}} -short


# Run the integration tests.
[group('test')]
int *FLAGS: check
  go test ./... -cover -vet=off -race {{FLAGS}} -run Integration

# HTML report for unit (default), int, e2e, or all tests.
[group('test')]
cover test='unit': check
  go test ./... -vet=off -coverprofile={{coverage_file}} \
  {{ if test == 'all' { '' } \
    else if test == 'int' { '-run Integration' } \
    else { '-short' } }}
  go tool cover -html={{coverage_file}}

# List the outdated direct dependencies (slow to run).
[group('dependencies')]
outdated:
  # (requires https://github.com/psampaz/go-mod-outdated).
  go list -u -m -json all | go-mod-outdated -update -direct

# Run go mod tidy and verify.
[group('dependencies')]
tidy:
  go mod tidy
  go mod verify

