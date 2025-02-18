.DEFAULT_GOAL := build

.PHONY:fmt vet build

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	tailwindcss -i ./ui/style/lovecms.css -o ./ui/static/lovecms.css
	go build -o ./dist/web ./cmd/web