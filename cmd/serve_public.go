package cmd

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/server/public"

	"github.com/spf13/cobra"
)

// servePublicCmd represents the public command
var servePublicCmd = &cobra.Command{
	Use:   "public",
	Short: "starts the public server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("public called")
	},
}

func init() {
	serveCmd.AddCommand(servePublicCmd)
}

func servePublic(ctx context.Context, options public.Options) error {
	server, err := public.NewServer(options)
	if err != nil {
		return err
	}
	server.Start(ctx)
	return nil
}
