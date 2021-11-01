package cmd

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/pkg/server/login"
	"github.com/spf13/cobra"
)

// serveLoginCmd represents the login command
var serveLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "starts the login server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("login called")
	},
}

func init() {
	serveCmd.AddCommand(serveLoginCmd)
}

func serveLogin(ctx context.Context, options login.Options) error {
	server, err := login.NewServer(options)
	if err != nil {
		return err
	}
	server.Start(ctx)
	return nil
}
