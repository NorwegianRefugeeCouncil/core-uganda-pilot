package main

import (
	"context"
	"flag"
	"github.com/nrc-no/core/api/cmd/server/app"
	serveroptions "github.com/nrc-no/core/api/pkg/server/options"
	"github.com/sirupsen/logrus"
	"k8s.io/apiserver/pkg/server"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	stopCh := server.SetupSignalHandler()
	go func() {
		<-stopCh
		cancel()
	}()

	options := &serveroptions.Options{}
	cmd := app.NewStartCoreServer(options, ctx)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	if err := cmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
