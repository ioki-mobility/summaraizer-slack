# Stage 1: Build the binary
FROM golang:1.22.2 AS builder

# Update ca-certificates
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Copy the Go project into the container
COPY . .

# Build the Go binary with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -a -o summaraizer-slack-server cmd/summaraizer-slack/main.go

# Stage 2: Create a minimal image with the compiled binary
FROM scratch

# Copy the CA certificates from the builder stage
# Otherwise the Go code is unable to make HTTP requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Set environment variables for the application
ENV PORT=8080

# Copy the binary from the builder stage
COPY --from=builder /app/summaraizer-slack-server /summaraizer-slack-server

# Set the entrypoint to the compiled binary
ENTRYPOINT ["/summaraizer-slack-server"]

# Pass runtime arguments to the binary
CMD ["/summaraizer-slack-server", "-port", "${PORT}"]
