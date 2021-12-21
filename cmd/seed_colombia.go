package cmd

import (
	"context"
	"net/url"

	"github.com/nrc-no/core/pkg/client"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/rest"
	"github.com/nrc-no/core/pkg/seeder"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var seedColombia = &cobra.Command{
	Use:   "colombia",
	Short: "Seed the database with default forms and fields for the Colombia context",
	RunE: func(cmd *cobra.Command, args []string) error {
		logging.SetLogLevel(zap.DebugLevel)

		ctx := context.Background()
		url, err := url.Parse(endpoint)
		if err != nil {
			return err
		}

		client := client.NewClientFromConfig(rest.Config{
			Scheme: url.Scheme,
			Host:   url.Host,
		})

		s, err := seeder.NewSeed(ctx, client)
		if err != nil {
			return err
		}
		return s.Seed(ctx, client, seeder.ColombiaContext)
	},
}

func init() {
  seedCmd.AddCommand(seedColombia)
}
