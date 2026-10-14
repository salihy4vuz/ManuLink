FROM docker.io/library/golang:alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main-api-binary main.go

# Run stage
FROM gcr.io/distroless/static-debian12

WORKDIR /app

# Copy the pre-built binary
COPY --from=builder /app/main-api-binary /app/main-api

# Command to run the executable
ENTRYPOINT ["/app/main-api"]

# Expose port 3000 to the outside world
EXPOSE 3000
