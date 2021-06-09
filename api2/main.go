package main

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/server"
)

func main() {
	ctx := context.Background()
	server.NewServer(ctx)
}
