package webapp

import (
	"html/template"
)

type Handler struct {
	template *template.Template
}

func NewHandler() (*Handler, error) {

	t, err := template.ParseGlob("pkg/webapp/templates/*.gohtml")
	if err != nil {
		return nil, err
	}
	h := &Handler{
		template: t,
	}
	return h, nil
}
