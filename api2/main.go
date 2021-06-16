package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/nrc-no/core-kafka/pkg/server"
	"log"
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

	// load .env file
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	srv.New(ctx)

	<-ctx.Done()
}
