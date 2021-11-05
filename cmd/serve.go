package cmd

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/options"
	"github.com/nrc-no/core/pkg/store"
	"github.com/spf13/viper"
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
		return storeSetup()
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

func storeSetup() error {
	logrus.SetFormatter(&logrus.TextFormatter{})
	setupSignal()
	if err := viper.Unmarshal(&coreOptions); err != nil {
		return err
	}

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
}
