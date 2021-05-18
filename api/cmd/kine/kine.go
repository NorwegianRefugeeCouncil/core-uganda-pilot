package main

import (
	"context"
	"fmt"
	"github.com/k3s-io/kine/pkg/endpoint"
)

func main() {
	ctx := context.Background()

	cfg, err := endpoint.Listen(ctx, endpoint.Config{
		Listener: "tcp://127.0.0.1:2379",
	})
	if err != nil {
		panic(err)
	}

	for _, endpoint := range cfg.Endpoints {
		fmt.Printf(endpoint)
	}

	<-ctx.Done()
}
