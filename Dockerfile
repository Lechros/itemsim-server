# 1. Build rure-go
FROM rust AS rust-builder

WORKDIR /app
RUN git clone --depth 1 https://github.com/rust-lang/regex
RUN cargo build --release --manifest-path ./regex/regex-capi/Cargo.toml

# 2. Build go binary
FROM golang:1.23 AS go-builder

# copy rure-go library to link
COPY --from=rust-builder /app/regex/target/release/*.so /usr/lib

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY internal internal
RUN go build -o app cmd/main.go

# 3. Production image
FROM ubuntu AS final

# add curl for healthcheck
RUN --mount=type=cache,target=/var/cache/apt \
    --mount=type=cache,target=/var/lib/apt \
    apt-get update && apt-get install -y curl

# copy rure-go library
COPY --from=rust-builder /app/regex/target/release/*.so /usr/lib

WORKDIR /app
COPY --from=go-builder /app/app .
COPY resources resources

EXPOSE 1323
ENTRYPOINT ["./app"]

HEALTHCHECK --start-period=10s --start-interval=2s \
  CMD curl -f http://localhost:1323/health || exit 1