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

# Copy the built binary from the builder stage
COPY --from=builder /run-app /usr/local/bin/
CMD ["run-app"]
