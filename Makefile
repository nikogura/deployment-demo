BINARY   := deployment-demo
VERSION  ?= dev
THEME    ?= green
HEALTH   ?= ok
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS = \
  -X github.com/nikogura/deployment-demo/pkg/demo.Version=$(VERSION) \
  -X github.com/nikogura/deployment-demo/pkg/demo.Theme=$(THEME) \
  -X github.com/nikogura/deployment-demo/pkg/demo.Health=$(HEALTH) \
  -X github.com/nikogura/deployment-demo/pkg/demo.BuildTime=$(BUILD_TIME)

REGISTRY ?= ghcr.io/nikogura
IMAGE    := $(REGISTRY)/deployment-demo

.PHONY: build build-ui build-go test lint clean docker-build docker-push help \
        v1.0.0 v1.0.2 v1.0.3 v1.0.4

help:
	@echo "Targets:"
	@echo "  build         Build UI + Go binary"
	@echo "  build-ui      Build the Next.js frontend"
	@echo "  build-go      Build the Go binary"
	@echo "  lint          Run all linters"
	@echo "  test          Run Go tests"
	@echo "  docker-build  Build Docker image (set VERSION, THEME, HEALTH)"
	@echo "  docker-push   Push image to registry"
	@echo "  v1.0.0        Build + push v1.0.0 (green, healthy)"
	@echo "  v1.0.2        Build + push v1.0.2 (blue, healthy)"
	@echo "  v1.0.3        Build + push v1.0.3 (red, BROKEN)"
	@echo "  v1.0.4        Build + push v1.0.4 (green, healthy — the fix)"

build: build-ui build-go

build-ui:
	@echo "Building UI..."
	cd pkg/ui && npm ci && npm run build

build-go:
	@echo "Building $(BINARY) $(VERSION) (theme=$(THEME), health=$(HEALTH))..."
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) ./cmd/deployment-demo

test:
	go test -v -race -cover ./...

lint:
	@echo "Running namedreturns linter..."
	@for pkg in $$(go list ./... | grep -v node_modules); do \
		namedreturns $$pkg || exit 1; \
	done
	@echo "Running golangci-lint..."
	golangci-lint run --timeout=5m
	@if [ -d pkg/ui/node_modules ]; then \
		echo "Running ESLint..."; \
		cd pkg/ui && ./node_modules/.bin/eslint 'app/**/*.{ts,tsx}' 'components/**/*.{ts,tsx}' 'lib/**/*.{ts,tsx}' 'types/**/*.{ts,tsx}' --max-warnings 0; \
	fi

clean:
	rm -rf bin/ pkg/ui/dist pkg/ui/.next pkg/ui/node_modules

# ---------------------------------------------------------------------------
# Docker
# ---------------------------------------------------------------------------
docker-build:
	@echo "Building Docker image $(IMAGE):$(VERSION) (theme=$(THEME), health=$(HEALTH))..."
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg THEME=$(THEME) \
		--build-arg HEALTH=$(HEALTH) \
		-t $(IMAGE):$(VERSION) .

docker-push:
	docker push $(IMAGE):$(VERSION)

# ---------------------------------------------------------------------------
# Versioned release shortcuts
# ---------------------------------------------------------------------------
v1.0.0:
	$(MAKE) docker-build docker-push VERSION=1.0.0 THEME=green HEALTH=ok

v1.0.2:
	$(MAKE) docker-build docker-push VERSION=1.0.2 THEME=blue HEALTH=ok

v1.0.3:
	$(MAKE) docker-build docker-push VERSION=1.0.3 THEME=red HEALTH=broken

v1.0.4:
	$(MAKE) docker-build docker-push VERSION=1.0.4 THEME=green HEALTH=ok
