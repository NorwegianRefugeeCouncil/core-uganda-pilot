FROM golang:1.16 AS builder

ARG git_tag
ARG git_commit

WORKDIR /go/src/github.com/nrc-no/core

ENV GO111MODULE=on
ADD go.mod .
ADD go.sum .

RUN go mod download

ADD . .

RUN go mod verify
RUN CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -a -ldflags '-linkmode external -extldflags "-static"' -o main .

# Build a small image
FROM scratch

COPY --from=builder /go/src/github.com/nrc-no/core/main /core

# Command to run
ENTRYPOINT ["/core"]
