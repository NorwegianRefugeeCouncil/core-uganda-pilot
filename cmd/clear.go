package cmd

import (
	"github.com/nrc-no/core/pkg/logging"
	loginstore "github.com/nrc-no/core/pkg/server/login/store"
	"github.com/nrc-no/core/pkg/store"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// migrateCmd does database migrations
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clears the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logging.NewLogger(ctx)
		factory, err := store.NewFactory(coreOptions.DSN)
		if err != nil {
			l.Error("failed to get factory", zap.Error(err))
			return err
		}
		db, err := factory.Get()
		if err != nil {
			l.Error("failed to get db", zap.Error(err))
			return err
		}
		if err := store.Clear(ctx, db); err != nil {
			l.Error("failed to clear store db", zap.Error(err))
			return err
		}
		if err := loginstore.Clear(db); err != nil {
			l.Error("failed to clear login store db", zap.Error(err))
			return err
		}
		l.Info("successfully cleared database")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
