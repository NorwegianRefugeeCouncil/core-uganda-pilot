FROM golang:1.16 AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
WORKDIR cmd/core
RUN go build -o main .

# Build a small image
FROM scratch

COPY --from=builder /build/cmd/core/main /
COPY --from=builder /build/web/templates /static/webapp/templates
COPY --from=builder /build/web/static /static/webapp/static

# Command to run
ENTRYPOINT ["/main"]
