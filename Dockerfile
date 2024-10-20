ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the entire app source
COPY . .

# Build the application, specifying the path to the main.go file inside the cmd directory
RUN go build -v -o /run-app ./cmd

FROM debian:bookworm

# Install CA certificates to ensure TLS connections can be verified
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the built binary from the builder stage
COPY --from=builder /run-app /usr/local/bin/

# Set the binary as the entrypoint
CMD ["run-app"]
