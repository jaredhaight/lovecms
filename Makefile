APP_NAME ?= lovecms

.PHONY: vet
vet:
	go vet .\...

.PHONY: test
test:
	go test -race -v -timeout 30s .\...

.PHONY: tailwind-watch
tailwind-watch:
	tailwindcss -i .\static\css\input.css -o .\static\css\style.css --watch

.PHONY: tailwind-build
tailwind-build:
	tailwindcss -i .\static\css\input.css -o .\static\css\style.min.css --minify

.PHONY: templ-watch
templ-watch:
	templ generate --watch

.PHONY: templ-generate
templ-generate:
	templ generate
	
.PHONY: dev
dev:
	go build -o C:\temp .\cmd\main.go && air

.PHONY: build
build:
	make tailwind-build
	make templ-generate
	go build -o .\bin\$(APP_NAME).exe .\cmd\main.go
