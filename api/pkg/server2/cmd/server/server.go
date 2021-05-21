package main

import (
	"flag"
	"github.com/nrc-no/core/api/pkg/server2/cmd/server/app"
	serveroptions "github.com/nrc-no/core/api/pkg/server2/options"
	"github.com/sirupsen/logrus"
)

func main() {
	options := &serveroptions.Options{}
	cmd := app.NewStartCoreServer(options)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	if err := cmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
