package server

import (
	"context"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/generic/server"
)

func (c *CompletedOptions) CreateIAMServer(ctx context.Context, genericOptions *server.GenericServerOptions) (*iam.Server, error) {
	iamServer, err := iam.NewServer(ctx, genericOptions)
	if err != nil {
		return nil, err
	}
	if err := iamServer.Init(ctx); err != nil {
		return nil, err
	}
	return iamServer, nil
}
