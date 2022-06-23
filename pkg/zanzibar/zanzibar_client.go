package zanzibar

import (
	"context"
	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"log"
)

type ZanzibarClient interface {
	WriteDB2UserRel(ctx context.Context, databaseId string, userId string) (*pb.WriteRelationshipsResponse, error)
}

type ZanzibarClientConfig struct {
	Token  string
	Prefix string
}

func NewZanzibarClient(c ZanzibarClientConfig) *zanzibarClient {
	client, err := authzed.NewClient(
		"grpc.authzed.com:443",
		grpcutil.WithBearerToken(c.Token),  // TODO get token from config
		grpcutil.WithSystemCerts(grpcutil.VerifyCA),
	)
	if err != nil {
		log.Fatalf("unable to initialize client: %s", err)
	}

	return &zanzibarClient{
		z: client,
		prefix: c.Prefix,
	}
}

type zanzibarClient struct {
	z      *authzed.Client
	prefix string
}

func (c *zanzibarClient) WriteDB2UserRel(ctx context.Context, databaseId string, userId string) (*pb.WriteRelationshipsResponse, error) {
	r := &pb.WriteRelationshipsRequest{
		Updates: []*pb.RelationshipUpdate{
			{
				Relationship: &pb.Relationship{
					Relation: "creator",
					Resource: &pb.ObjectReference{
						ObjectType: c.prefix + "/database", // TODO get prefix from config
						ObjectId:   databaseId,
					},
					Subject: &pb.SubjectReference{
						Object: &pb.ObjectReference{
							ObjectType: c.prefix + "/user", // TODO get prefix from config
							ObjectId:   userId,
						},
					},
				},
				Operation: pb.RelationshipUpdate_OPERATION_CREATE,
			},
		},
	}

	resp, err := c.z.WriteRelationships(ctx, r)

	if err != nil {
		log.Fatalf("failed to create relationship between database and creator: %s, %s", err, resp)
		return nil, err
	}
	return resp, nil
}
