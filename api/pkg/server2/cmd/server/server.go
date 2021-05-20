package main

import (
	"flag"
	"github.com/nrc-no/core/api/pkg/server2/cmd/server/app"
	serveroptions "github.com/nrc-no/core/api/pkg/server2/options"
	"github.com/sirupsen/logrus"
	"k8s.io/apiserver/pkg/server"
)

func main() {
	stopCh := server.SetupSignalHandler()
	options := &serveroptions.Options{}
	cmd := app.NewStartCoreServer(options, stopCh)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	if err := cmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
