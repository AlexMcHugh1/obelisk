SHELL := /bin/bash

APP_NAME ?= obelisk
IMAGE_NAME ?= $(APP_NAME)
IMAGE_TAG ?= latest
IMAGE_REF ?= $(IMAGE_NAME):$(IMAGE_TAG)

DOCKERFILE ?= Dockerfile
COMPOSE_FILE ?= docker-compose.yml

GO ?= go
PKGS := ./...
CMD_PATH := ./cmd/server

GOLANGCI_LINT ?= golangci-lint
STATICCHECK ?= staticcheck
GOSEC ?= gosec
TRIVY ?= trivy
GRYPE ?= grype
DOCKER ?= docker

.PHONY: help fmt test vet lint staticcheck scan security build build-image compose-up compose-down compose-logs trivy-scan grype-scan clean check-tools

help:
	@echo "Available targets:"
	@echo "  fmt           - Format Go code"
	@echo "  test          - Run tests"
	@echo "  vet           - Run go vet"
	@echo "  lint          - Run golangci-lint"
	@echo "  staticcheck   - Run staticcheck"
	@echo "  scan          - Run gosec security scan"
	@echo "  security      - Run lint + staticcheck + gosec"
	@echo "  build         - Build Go binary"
	@echo "  build-image   - Build Docker image"
	@echo "  compose-up    - Start docker compose stack"
	@echo "  compose-down  - Stop docker compose stack"
	@echo "  compose-logs  - Tail docker compose logs"
	@echo "  trivy-scan    - Scan filesystem and image with Trivy"
	@echo "  grype-scan    - Scan image with Grype"
	@echo "  clean         - Remove build artifacts"

fmt:
	$(GO) fmt $(PKGS)

test:
	$(GO) test -race -cover $(PKGS)

vet:
	$(GO) vet $(PKGS)

lint:
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || { echo "golangci-lint is not installed"; exit 1; }
	$(GOLANGCI_LINT) run ./...

staticcheck:
	@command -v $(STATICCHECK) >/dev/null 2>&1 || { echo "staticcheck is not installed"; exit 1; }
	$(STATICCHECK) $(PKGS)

scan:
	@command -v $(GOSEC) >/dev/null 2>&1 || { echo "gosec is not installed"; exit 1; }
	$(GOSEC) ./...

security: lint staticcheck scan

build:
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -o bin/server $(CMD_PATH)

build-image:
	$(DOCKER) build -f $(DOCKERFILE) -t $(IMAGE_REF) .

compose-up:
	$(DOCKER) compose -f $(COMPOSE_FILE) up -d --build

compose-down:
	$(DOCKER) compose -f $(COMPOSE_FILE) down

compose-logs:
	$(DOCKER) compose -f $(COMPOSE_FILE) logs -f

trivy-scan:
	@command -v $(TRIVY) >/dev/null 2>&1 || { echo "trivy is not installed"; exit 1; }
	$(TRIVY) fs --scanners vuln,secret,misconfig .
	$(TRIVY) image --scanners vuln,secret,misconfig $(IMAGE_REF)

grype-scan:
	@command -v $(GRYPE) >/dev/null 2>&1 || { echo "grype is not installed"; exit 1; }
	$(GRYPE) $(IMAGE_REF)

clean:
	rm -rf bin