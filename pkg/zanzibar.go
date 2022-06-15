package main

import (
	"context"
	"log"

	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
)

const schema = `definition nrc_test/user {}

definition nrc_test/post {
	relation reader: nrc_test/user
	relation writer: nrc_test/user

	permission read = reader + writer
	permission write = writer
}`

func main() {
	client, err := authzed.NewClient(
		"grpc.authzed.com:443",
		grpcutil.WithBearerToken("tc_nrc_test_default_token"),
		grpcutil.WithSystemCerts(grpcutil.VerifyCA),
	)
	if err != nil {
		log.Fatalf("unable to initialize client: %s", err)
	}

	request := &pb.WriteSchemaRequest{Schema: schema}
	resp, err := client.WriteSchema(context.Background(), request)
	if err != nil {
		log.Fatalf("failed to write schema: %s", resp)
	}
}
