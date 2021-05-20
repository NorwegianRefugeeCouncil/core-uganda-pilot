module github.com/nrc-no/core/api

go 1.16

require (
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/emicklei/go-restful v2.9.5+incompatible
	github.com/emicklei/go-restful/v3 v3.5.1 // indirect
	github.com/go-openapi/spec v0.19.5 // indirect
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/google/gofuzz v1.1.0
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/k3s-io/kine v0.6.0
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	go.mongodb.org/mongo-driver v1.5.2
	go.uber.org/multierr v1.5.0 // indirect
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	k8s.io/api v0.21.0
	k8s.io/apiextensions-apiserver v0.21.0
	k8s.io/apimachinery v0.21.0
	k8s.io/apiserver v0.21.0
	k8s.io/client-go v0.21.0
	k8s.io/code-generator v0.21.0
	k8s.io/klog v1.0.0
	k8s.io/klog/v2 v2.8.0
	k8s.io/kube-openapi v0.0.0-20210421082810-95288971da7e
	sigs.k8s.io/structured-merge-diff/v4 v4.1.1
)

replace (
	k8s.io/api => ./third-party/kubernetes/staging/src/k8s.io/api
	k8s.io/apiextensions-apiserver => ./third-party/kubernetes/staging/src/k8s.io/apiextensions-apiserver
	k8s.io/apimachinery => ./third-party/kubernetes/staging/src/k8s.io/apimachinery
	k8s.io/apiserver => ./third-party/kubernetes/staging/src/k8s.io/apiserver
	k8s.io/client-go => ./third-party/kubernetes/staging/src/k8s.io/client-go
	k8s.io/code-generator => ./third-party/kubernetes/staging/src/k8s.io/code-generator
	k8s.io/component-base => ./third-party/kubernetes/staging/src/k8s.io/component-base
)
