package client

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

func NewZanzibarClient() *zanzibarClient {
	client, err := authzed.NewClient(
		"grpc.authzed.com:443",
		grpcutil.WithBearerToken("tc_nrctest_default_token_e573a4e880dbb1219af014722b4e70802fdc53b01508ac19fc8e2b2f50edf94ef477c21779ebc823787c1809c30fc84e69d177dd424cda2738ffaafae8b9e8e4"),  // TODO get token from config
		grpcutil.WithSystemCerts(grpcutil.VerifyCA),
	)
	if err != nil {
		log.Fatalf("unable to initialize client: %s", err)
	}

	return &zanzibarClient{
		z: client,
	}

}

type zanzibarClient struct {
	z *authzed.Client
}

func (c *zanzibarClient) WriteDB2UserRel(ctx context.Context, databaseId string, userId string) (*pb.WriteRelationshipsResponse, error) {
	r := &pb.WriteRelationshipsRequest{
		Updates: []*pb.RelationshipUpdate{
			{
				Relationship: &pb.Relationship{
					Relation: "creator",
					Resource: &pb.ObjectReference{
						ObjectType: "nrctest/database", // TODO get prefix from config
						ObjectId:   databaseId,
					},
					Subject: &pb.SubjectReference{
						Object: &pb.ObjectReference{
							ObjectType: "nrctest/user", // TODO get prefix from config
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
