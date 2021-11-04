package cmd

import (
	"context"
	"github.com/nrc-no/core/pkg/server/options"
	"github.com/nrc-no/core/pkg/store"
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
		logrus.SetFormatter(&logrus.TextFormatter{})
		setupSignal()
		if err := viper.Unmarshal(&coreOptions); err != nil {
			return err
		}

		if len(coreOptions.Log.Level) > 0 {
			logLevel, err := logrus.ParseLevel(coreOptions.Log.Level)
			if err != nil {
				return err
			}
			logrus.SetLevel(logLevel)
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
