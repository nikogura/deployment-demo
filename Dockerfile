# ---------- Stage 1: Build the UI ----------
FROM node:22-alpine AS ui-builder
WORKDIR /app/pkg/ui
COPY pkg/ui/package.json pkg/ui/package-lock.json* ./
RUN npm ci
COPY pkg/ui/ .
RUN npm run build

# ---------- Stage 2: Build the Go binary ----------
FROM golang:1.26-alpine AS go-builder

ARG VERSION=dev
ARG THEME=green
ARG HEALTH=ok

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=ui-builder /app/pkg/ui/dist pkg/ui/dist

RUN CGO_ENABLED=0 go build \
  -ldflags "\
    -X github.com/nikogura/deployment-demo/pkg/demo.Version=${VERSION} \
    -X github.com/nikogura/deployment-demo/pkg/demo.Theme=${THEME} \
    -X github.com/nikogura/deployment-demo/pkg/demo.Health=${HEALTH} \
    -X github.com/nikogura/deployment-demo/pkg/demo.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -o /deployment-demo ./cmd/deployment-demo

# ---------- Stage 3: Distroless runtime ----------
FROM gcr.io/distroless/static:nonroot
COPY --from=go-builder /deployment-demo /deployment-demo
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/deployment-demo"]
