package handlers

import "net/http"

type Identificator interface {
	Identify(req *http.Request) string
}
