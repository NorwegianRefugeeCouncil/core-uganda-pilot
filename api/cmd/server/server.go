package main

import (
	"flag"
	app2 "github.com/nrc-no/core/api/cmd/server/app"
	serveroptions "github.com/nrc-no/core/api/pkg/server2/options"
	"github.com/sirupsen/logrus"
)

func main() {
	options := &serveroptions.Options{}
	cmd := app2.NewStartCoreServer(options)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	if err := cmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
