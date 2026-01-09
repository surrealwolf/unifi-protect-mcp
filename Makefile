.PHONY: help build run test fmt lint check docker-build docker-run docker-login docker-pull-base docker-push clean

help:
	@echo "Unifi Protect MCP Server - Make Commands"
	@echo "========================================"
	@echo "  make build          - Build the protect MCP binary"
	@echo "  make run            - Run the protect MCP server"
	@echo "  make test           - Run tests"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Run linter"
	@echo "  make check          - Run all checks (fmt, lint, test)"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-run     - Run Docker container"
	@echo "  make docker-login   - Login to Harbor registry"
	@echo "  make docker-pull-base - Pull base images from Harbor cache"
	@echo "  make docker-push    - Build and push to Harbor"
	@echo "  make clean          - Clean build artifacts"

build:
	go build -o bin/unifi-protect-mcp ./cmd

run: build
	./bin/unifi-protect-mcp

test:
	go test -v -cover ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

check: fmt lint test

# Harbor registry configuration
HARBOR_REGISTRY ?= harbor.dataknife.net
HARBOR_PROJECT ?= library
IMAGE_NAME ?= unifi-protect-mcp
IMAGE_TAG ?= latest
FULL_IMAGE = $(HARBOR_REGISTRY)/$(HARBOR_PROJECT)/$(IMAGE_NAME):$(IMAGE_TAG)

docker-login:
	@echo "Logging in to Harbor registry..."
	@docker login $(HARBOR_REGISTRY) \
		-u '$${HARBOR_USERNAME:-robot$$library+ci-builder}' \
		-p '$${HARBOR_PASSWORD}'

docker-pull-base:
	@echo "Pulling base images from Harbor cache..."
	@docker pull $(HARBOR_REGISTRY)/dockerhub/library/golang:1.23-alpine || true
	@docker pull $(HARBOR_REGISTRY)/dockerhub/library/alpine:3.18 || true

docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(FULL_IMAGE)

docker-build-harbor: docker-pull-base
	docker build \
		--cache-from $(FULL_IMAGE) \
		-t $(IMAGE_NAME):$(IMAGE_TAG) \
		-t $(FULL_IMAGE) \
		.

docker-push: docker-build-harbor docker-login
	docker push $(FULL_IMAGE)

docker-run: docker-build
	docker run --rm \
		-e UNIFI_PROTECT_URL="$${UNIFI_PROTECT_URL}" \
		-e UNIFI_PROTECT_USERNAME="$${UNIFI_PROTECT_USERNAME}" \
		-e UNIFI_PROTECT_PASSWORD="$${UNIFI_PROTECT_PASSWORD}" \
		$(FULL_IMAGE)

clean:
	rm -rf bin/
	go clean
