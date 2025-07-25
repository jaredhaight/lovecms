APP_NAME ?= lovecms

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test -race -v -timeout 30s ./...

.PHONY: test-coverage
test-coverage:
	go test -race -v -timeout 30s -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-short
test-short:
	go test -short -race -v ./...

.PHONY: benchmark
benchmark:
	go test -bench=. -benchmem ./...

.PHONY: tailwind-watch
tailwind-watch:
	npx @tailwindcss/cli -i ./static/css/input.css -o ./static/css/lovecms.css --watch

.PHONY: tailwind-build
tailwind-build:
	npx @tailwindcss/cli -i ./static/css/input.css -o ./static/css/lovecms.css

.PHONY: templ-watch
templ-watch:
	templ generate --watch

.PHONY: templ-generate
templ-generate:
	templ generate

.PHONY: build
build:
	make tailwind-build
	make templ-generate
	go build -o ./bin/$(APP_NAME) ./cmd/main.go
