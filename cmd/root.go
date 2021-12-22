package cmd

import (
	"github.com/drone/envsubst"
	"github.com/mitchellh/mapstructure"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/options"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"reflect"
)

var cfgFiles []string
var v *viper.Viper

var decodeHook = viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
	func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() == reflect.String {
			return envsubst.EvalEnv(data.(string))
		}
		return data, nil
	},
	mapstructure.StringToTimeDurationHookFunc(),
	mapstructure.StringToSliceHookFunc(","),
))

func unmarshalConfig(out *options.Options) error {
	return v.Unmarshal(out, decodeHook)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "core",
	Short: "Data collection and case management for humanitarian organizations",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return unmarshalConfig(&coreOptions)
	},
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
	v = viper.New()
	v.AutomaticEnv() // read in environment variables that match
	for i, file := range cfgFiles {
		v.SetConfigFile(file)
		var err error
		if i == 0 {
			err = v.ReadInConfig()
		} else {
			err = v.MergeInConfig()
		}
		if err == nil {
			l.Info("using config file", zap.String("config_file", v.ConfigFileUsed()))
		} else {
			l.Error("failed to use config file", zap.Error(err))
			panic(err)
		}
	}

	if len(cfgFiles) > 0 {
		v.WatchConfig()
		l.Info("watch started")
	}

}
