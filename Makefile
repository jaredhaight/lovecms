APP_NAME ?= lovecms

.PHONY: vet
vet:
	go vet .\...

.PHONY: test
test:
	go test -race -v -timeout 30s .\...

.PHONY: tailwind-watch
tailwind-watch:
	tailwindcss -i .\static\css\input.css -o .\static\css\lovecms.css --watch

.PHONY: tailwind-build
tailwind-build:
	tailwindcss -i .\static\css\input.css -o .\static\css\lovecms.css

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
	go build -o .\bin\$(APP_NAME).exe .\cmd\main.go
