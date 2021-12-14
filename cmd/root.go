package cmd

import (
	"github.com/nrc-no/core/pkg/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var cfgFiles []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "core",
	Short: "Data collection and case management for humanitarian organizations",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringSliceVar(&cfgFiles, "config", []string{}, "config files")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	l := logging.NewLogger(ctx)

	viper.AutomaticEnv() // read in environment variables that match

	for i, file := range cfgFiles {
		viper.SetConfigFile(file)
		var err error
		if i == 0 {
			err = viper.ReadInConfig()
		} else {
			err = viper.MergeInConfig()
		}
		if err == nil {
			l.Info("using config file", zap.String("config_file", viper.ConfigFileUsed()))
		} else {
			l.Error("failed to use config file", zap.Error(err))
			panic(err)
		}
	}

	l.Info("finished init config. Start watch")

	if len(cfgFiles) > 0 {
		viper.WatchConfig()
	}

	l.Info("watch started")

}
