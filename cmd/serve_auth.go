package cmd

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/server/auth"
	"github.com/spf13/cobra"
)

// serveAuthCmd represents the admin command
var serveAuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "starts the auth server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("auth called")
	},
}

func init() {
	serveCmd.AddCommand(serveAdminCmd)
}

func serveAuth(ctx context.Context, options auth.Options) error {
	server, err := auth.NewServer(options)
	if err != nil {
		return err
	}
	server.Start(ctx)
	return nil
}
