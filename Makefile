# Variables
GO = go
FMT = gofmt

build:
	$(GO) build

fmt:
	$(FMT) -w .

.PHONY: build fmt
