package main

import (
	"github.com/nrc-no/core/api/pkg/server2"
	"os"
	"os/signal"
)

func main() {

	options := server2.Config{
		ListenAddress: ":8888",
		StorageConfig: &server2.MongoConfig{
			Address: "mongodb://root:pass12345@localhost:27017",
		},
	}

	srv := options.Complete().New()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	stopCh := make(chan struct{}, 1)
	go func() {
		<-signalCh
		stopCh <- struct{}{}
	}()

	if err := srv.Run(stopCh); err != nil {
		panic(err)
	}

}
