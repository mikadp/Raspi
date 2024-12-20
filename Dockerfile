FROM golang:alpine AS builder

# Setting the working directory in the container
WORKDIR /app

# Install build tools only in the build stage
 RUN apt-get update 
 RUN apt-get update && apt-get install -y gcc g++ libc6-dev
 

# Copy the go.mod and go.sum files to the container first
COPY go.mod ./ go.sum ./

# Download Go modules
RUN go mod download

# Copy all files to the container
COPY . .

# Build the Go app
RUN GOARCH=amd64 go build -o raspigoapp .
ENV CGO_ENABLED=0  
#GOOS=linux 

# minimal base image to keep the final image small
FROM alpine:latest

# Working directory inside the new image
WORKDIR /root/

# Install runtime dependencies
RUN apk add --no-cache libpq

# Copying the compiled binary from the build container
COPY --from=builder /app/raspigoapp .

# executable permissions for the binary
RUN chmod +x raspigoapp

# run the app
CMD ["./raspigoapp"]