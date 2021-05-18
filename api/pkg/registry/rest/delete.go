package rest

import "github.com/nrc-no/core/apps/api/pkg/runtime"

type RESTDeleteStrategy interface {
	runtime.ObjectTyper
}
