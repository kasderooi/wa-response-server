ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm AS builder

# Set working directory inside the container
WORKDIR /usr/src/app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the entire project into the container
COPY . .

# Build the Go application from /cmd/wa-server
RUN go build -v -o /run-app ./cmd/wa-server


# Create a smaller final image
FROM debian:bookworm

# Copy the built binary from the builder stage
COPY --from=builder /run-app /usr/local/bin/

# Run the built app
CMD ["run-app"]
