APP_NAME ?= lovecms

.PHONY: vet
vet:
	go vet .\...

.PHONY: test
test:
	go test -race -v -timeout 30s .\...

.PHONY: tailwind-watch
tailwind-watch:
	tailwindcss -i .\static\css\input.css -o .\static\css\lovecms.min.css --watch

.PHONY: tailwind-build
tailwind-build:
	tailwindcss -i .\static\css\input.css -o .\static\css\lovecms.min.css --minify

.PHONY: dev
dev:
	go build -o C:\temp .\cmd\main.go && air

.PHONY: build
build:
	make tailwind-build
	go build -o .\bin\$(APP_NAME).exe .\cmd\main.go
