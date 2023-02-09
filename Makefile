help:
	@echo "You can perform the following:"
	@echo ""
	@echo "  check         Format, vet, and test Go code"
	@echo "  cover         Show test coverage in html"
	@echo "  lint          Lint Go code using staticcheck"

check:
	@echo 'Formatting, vetting, and testing Go code'
	go fmt ./...
	go vet ./...
	go test ./... -cover

cover:
	@echo 'Test coverage in html'
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

lint:
	@echo 'Linting code using staticcheck'
	staticcheck -f stylish ./...
