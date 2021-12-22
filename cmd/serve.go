package cmd

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/options"
	"github.com/nrc-no/core/pkg/store"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var coreOptions options.Options
var osSignal = make(chan os.Signal)
var doneSignal = make(chan struct{})
var ctx = context.Background()
var factory store.Factory

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "base command for starting servers",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		setupSignal()
		configCtx := context.Background()
		configLog := logging.NewLogger(configCtx)
		v.OnConfigChange(func(in fsnotify.Event) {
			var changedConfig options.Options
			configLog.Info("detected configuration change")
			if err := unmarshalConfig(&changedConfig); err != nil {
				configLog.Error("failed to unmarshal on config change", zap.Error(err))
				return
			}
			coreOptions = changedConfig
		})
		switch coreOptions.Log.Level {
		case "debug":
			logging.SetLogLevel(zapcore.DebugLevel)
		case "info":
			logging.SetLogLevel(zapcore.InfoLevel)
		case "warn":
			logging.SetLogLevel(zapcore.WarnLevel)
		case "error":
			logging.SetLogLevel(zapcore.ErrorLevel)
		case "dpanic":
			logging.SetLogLevel(zapcore.DPanicLevel)
		case "panic":
			logging.SetLogLevel(zapcore.PanicLevel)
		case "fatal":
			logging.SetLogLevel(zapcore.FatalLevel)
		}

		return nil
	},
}

var factoryLock sync.Mutex

func initStoreFactory() error {
	factoryLock.Lock()
	defer factoryLock.Unlock()
	if factory != nil {
		return nil
	}
	f, err := store.NewDynamicFactory(func() (string, error) {
		return coreOptions.DSN, nil
	})
	if err != nil {
		return err
	}
	factory = f
	return nil
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
