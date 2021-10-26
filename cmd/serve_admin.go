package cmd

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/server/admin"
	"github.com/spf13/cobra"
)

// serveAdminCmd represents the admin command
var serveAdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "starts the admin server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("admin called")
	},
}

func init() {
	serveCmd.AddCommand(serveAdminCmd)
}

func serveAdmin(ctx context.Context, options admin.Options) error {
	server, err := admin.NewServer(options)
	if err != nil {
		return err
	}
	server.Start(ctx)
	return nil
}
