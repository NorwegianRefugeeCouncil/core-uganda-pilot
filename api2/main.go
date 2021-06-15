package main

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/server"
)

func main() {
	ctx := context.Background()
	serverOptions := server.Options{
		MongoUsername: "root",
		MongoPassword: "example",
	}
	srv, err := serverOptions.Complete(ctx)
	if err != nil {
		panic(err)
	}
	srv.New(ctx)

	<-ctx.Done()
}
