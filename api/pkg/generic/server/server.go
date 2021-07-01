package server

import "github.com/spf13/pflag"

type GenericServerOptions struct {
	ListenAddress string
	MongoHosts    []string
	MongoDatabase string
	MongoUsername string
	MongoPassword string
	Environment   string
}

func NewServerOptions() *GenericServerOptions {
	return &GenericServerOptions{
		ListenAddress: ":9001",
		MongoHosts:    []string{"mongo://localhost:27017"},
	}
}

func (o *GenericServerOptions) WithMongoHosts(hosts []string) *GenericServerOptions {
	o.MongoHosts = hosts
	return o
}
func (o *GenericServerOptions) WithMongoDatabase(mongoDatabase string) *GenericServerOptions {
	o.MongoDatabase = mongoDatabase
	return o
}
func (o *GenericServerOptions) WithMongoUsername(mongoUsername string) *GenericServerOptions {
	o.MongoUsername = mongoUsername
	return o
}
func (o *GenericServerOptions) WithMongoPassword(mongoPassword string) *GenericServerOptions {
	o.MongoPassword = mongoPassword
	return o
}
func (o *GenericServerOptions) WithListenAddress(address string) *GenericServerOptions {
	o.ListenAddress = address
	return o
}
func (o *GenericServerOptions) WithEnvironment(environment string) *GenericServerOptions {
	o.Environment = environment
	return o
}

func (o *GenericServerOptions) Flags(fs pflag.FlagSet) {
	fs.StringVar(&o.ListenAddress, "listen-address", o.ListenAddress, "Server listen address")
	fs.StringSliceVar(&o.MongoHosts, "mongo-url", o.MongoHosts, "Mongo url")
	fs.StringVar(&o.MongoDatabase, "mongo-database", o.MongoDatabase, "Mongo database")
	fs.StringVar(&o.MongoUsername, "mongo-username", o.MongoUsername, "Mongo username")
	fs.StringVar(&o.MongoPassword, "mongo-password", o.MongoPassword, "Mongo password")
}
