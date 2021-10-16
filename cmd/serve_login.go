package cmd

import (
	"fmt"

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
