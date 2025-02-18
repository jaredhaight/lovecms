.DEFAULT_GOAL := build

.PHONY:fmt vet build

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	cd ./ui && npx @tailwindcss/cli -i .\style\lovecms.css -o .\static\lovecms.css
	go build ./cmd/web