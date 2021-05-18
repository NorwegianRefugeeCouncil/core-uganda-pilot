package main

import (
	"flag"
	"github.com/nrc-no/core/api/cmd/server/app"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/klog/v2"
)

func main() {
	stopCh := genericapiserver.SetupSignalHandler()
	options := app.NewCoreServerOptions()
	cmd := app.NewCommandStartCoreServer(options, stopCh)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	if err := cmd.Execute(); err != nil {
		klog.Fatal(err)
	}
}
