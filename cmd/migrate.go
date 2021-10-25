package cmd

import (
	"github.com/nrc-no/core/pkg/bla/store"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// migrateCmd does database migrations
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Executes migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.Unmarshal(&coreOptions); err != nil {
			return err
		}
		factory, err := store.NewFactory(coreOptions.DSN)
		if err != nil {
			return err
		}
		db, err := factory.Get()
		if err != nil {
			return err
		}
		if err := store.Migrate(db); err != nil {
			return err
		}
		logrus.Info("Successfully applied migrations")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
