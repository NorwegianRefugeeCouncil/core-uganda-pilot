package main

import (
	"context"
	"flag"
	"github.com/nrc-no/core/cmd/app"
	"github.com/nrc-no/core/pkg/apps/iam"
	"github.com/nrc-no/core/pkg/server"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

func main() {
	ctx := context.Background()
	options := server.NewOptions()
	cmd := app.LaunchCommand(ctx, options)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)

	// Generate TS types for client
	converter := typescriptify.New().
		Add(iam.Party{})
	err := converter.ConvertToFile("models.ts")
	if err != nil {
		panic(err.Error())
	}
	// end of TS type generation

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
	<-ctx.Done()
}
