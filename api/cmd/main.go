package main

import (
	"context"
	"flag"
	"github.com/nrc-no/core/cmd/app"
	"github.com/nrc-no/core/pkg/server"
)

func main() {
	ctx := context.Background()
	options := server.NewOptions()
	cmd := app.LaunchCommand(ctx, options)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
	<-ctx.Done()
}
