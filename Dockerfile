FROM golang:latest AS builder

# Setting the working directory in the container
WORKDIR /app

# Copy the go.mod and go.sum files to the container first
COPY go.mod ./ go.sum ./

# Download Go modules
RUN go mod download

# Copy all files to the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o raspigoapp

# minimal base image to keep the final image small
FROM alpine:latest

# Working directory inside the new image
WORKDIR /app

# Copying the compiled binary from the build container
COPY --from=builder /app/raspigoapp .

# executable permissions for the binary
RUN chmod +x raspigoapp

# run the app
CMD ["./raspigoapp"]