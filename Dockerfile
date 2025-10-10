# 1. Build go binary
FROM golang:1.23 AS go-builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY internal internal
RUN go build -o app cmd/main.go

# 2. Production image
FROM alpine AS final

# Install required packages
# gcompat: Go binary compatibility for Alpine Linux
RUN apk add --no-cache gcompat

WORKDIR /app
COPY --from=go-builder /app/app .
COPY resources resources

EXPOSE 1323

# Health check using wget (built into Alpine)
HEALTHCHECK --start-period=10s --interval=30s --timeout=3s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:1323/health || exit 1

ENTRYPOINT ["./app"]
