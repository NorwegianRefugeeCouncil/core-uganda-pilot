package cmd

import (
	"github.com/nrc-no/core/pkg/server/options"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestConfigReplacesEnvVar(t *testing.T) {
	testCmd := &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	rootCmd.AddCommand(testCmd)

	tests := []struct {
		name   string
		config string
		env    map[string]string
		assert func(t *testing.T, options options.Options)
	}{
		{
			name:   "simple",
			config: "dsn: prefix/${BLA}/suffix",
			env:    map[string]string{"BLA": "BLA_VALUE"},
			assert: func(t *testing.T, options options.Options) {
				assert.Equal(t, "prefix/BLA_VALUE/suffix", options.DSN)
			},
		}, {
			name:   "default value",
			config: "dsn: prefix/${BLA:-DEFAULT_BLA}/suffix",
			assert: func(t *testing.T, options options.Options) {
				assert.Equal(t, "prefix/DEFAULT_BLA/suffix", options.DSN)
			},
		}, {
			name: "bool replace",
			env: map[string]string{
				"SNIP": "true",
			},
			config: `
serve:
  forms_api:
    cors:
      enabled: ${SNIP}
`,
			assert: func(t *testing.T, options options.Options) {
				assert.True(t, options.Serve.FormsApi.Cors.Enabled)
			},
		}, {
			name: "int replace",
			env: map[string]string{
				"SNIP": "9000",
			},
			config: `
serve:
  forms_api:
    port: ${SNIP}
`,
			assert: func(t *testing.T, options options.Options) {
				assert.Equal(t, 9000, options.Serve.FormsApi.Port)
			},
		},
	}

	cfgFile, err := ioutil.TempFile("", "test-parse-config-*.yaml")
	if !assert.NoError(t, err) {
		return
	}
	defer cfgFile.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Resetting config files list
			cfgFiles = []string{}

			var envs []string
			defer func() {
				for _, env := range envs {
					os.Unsetenv(env)
				}
			}()

			for key, value := range test.env {
				if err := os.Setenv(key, value); !assert.NoError(t, err) {
					return
				}
				envs = append(envs, key)

			}

			args := []string{"test"}
			if len(test.config) != 0 {
				if _, err := cfgFile.Seek(0, io.SeekStart); !assert.NoError(t, err) {
					return
				}

				if err := cfgFile.Truncate(0); !assert.NoError(t, err) {
					return
				}

				if _, err := cfgFile.WriteString(test.config); !assert.NoError(t, err) {
					return
				}
				args = append(args, "--config", cfgFile.Name())
			}

			rootCmd.SetArgs(args)
			if err := rootCmd.Execute(); !assert.NoError(t, err) {
				return
			}

			test.assert(t, coreOptions)
		})
	}

}
