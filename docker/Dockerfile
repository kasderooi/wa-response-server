# Use official Go image
FROM golang:1.23.4

WORKDIR /app

# Copy and build the app
COPY . .
RUN go mod tidy
RUN go build -o cmd/main .

# Expose port 8080
EXPOSE 8080

# Run the backend
CMD ["/app/main"]
