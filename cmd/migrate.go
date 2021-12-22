package cmd

import (
	"github.com/nrc-no/core/pkg/logging"
	loginstore "github.com/nrc-no/core/pkg/server/login/store"
	"github.com/nrc-no/core/pkg/store"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// migrateCmd does database migrations
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Executes migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logging.NewLogger(ctx)
		factory, err := store.NewFactory(coreOptions.DSN)
		if err != nil {
			l.Error("failed to create factory", zap.Error(err))
			return err
		}
		db, err := factory.Get()
		if err != nil {
			l.Error("failed to get database connection", zap.Error(err))
			return err
		}
		if err := store.Migrate(db); err != nil {
			l.Error("failed to migrate store", zap.Error(err))
			return err
		}
		if err := loginstore.Migrate(db); err != nil {
			l.Error("failed to migrate login store", zap.Error(err))
			return err
		}
		l.Info("successfully applied migrations")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
