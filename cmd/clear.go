package cmd

import (
	loginstore "github.com/nrc-no/core/pkg/server/login/store"
	"github.com/nrc-no/core/pkg/store"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// migrateCmd does database migrations
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clears the database",
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
		if err := store.Clear(ctx, db); err != nil {
			return err
		}
		if err := loginstore.Clear(db); err != nil {
			return err
		}
		logrus.Info("Successfully cleared database!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
