package cmd

import (
	"context"
	"github.com/nrc-no/core/pkg/bla/options"
	"github.com/nrc-no/core/pkg/bla/store"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

var coreOptions options.Options
var osSignal = make(chan os.Signal)
var doneSignal = make(chan struct{})
var serveCtx = context.Background()
var factory store.Factory

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "base command for starting servers",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		setupSignal()
		if err := viper.Unmarshal(&coreOptions); err != nil {
			return err
		}
		var err error
		factory, err = store.NewFactory(coreOptions.DSN)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
func setupSignal() {
	go func() {
		<-osSignal
		doneSignal <- struct{}{}
	}()
}
