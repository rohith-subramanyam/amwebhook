# Dockerfile References: https://docs.docker.com/engine/reference/builder/.

# Start from the latest golang base image.
FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Add Maintainer Info.
#LABEL maintainer="Rohith Subramanyam <rohith.subramanyam@nutanix.com>"

# Set the Current Working Directory inside the container.
WORKDIR /app

# Copy go mod and sum files.
# All the dependent modules do not conform to go mod yet.
#COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and
# go.sum files are not changed.
# All the dependent modules do not conform to go mod yet.
#RUN go mod download
RUN go get github.com/prometheus/alertmanager/template

# Copy the source from the current directory to the Working Directory inside the
# container.
COPY . .

# Build the Go app.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o amwebhook .


######## Start a new stage from scratch ########
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage.
COPY --from=builder /app/amwebhook .

# Expose port 8080 to the outside world.
EXPOSE 8080

# Command to run the executable.
CMD ["./amwebhook"]
