GO ?= go
CV := ./bin/cv

.PHONY: setup fmt test vet build resume clean
setup:
	$(GO) mod download
	@command -v lualatex >/dev/null || { echo "Missing lualatex (install texlive-luatex and recommended fonts/packages)"; exit 1; }
	@command -v pdftotext >/dev/null || { echo "Missing pdftotext (install poppler-utils)"; exit 1; }
fmt:
	$(GO) fmt ./...
test:
	$(GO) test ./...
vet:
	$(GO) vet ./...
build:
	mkdir -p bin
	$(GO) build -o $(CV) ./cmd/cv
resume: build
	$(CV) build --lang fr --target senior-backend
	$(CV) build --lang en --target senior-backend
clean:
	rm -rf bin output/*.aux output/*.log output/*.out output/*.tex output/*.pdf
