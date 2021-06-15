package webapp

import (
	"fmt"
	"html/template"
	"os"
)

type Handler struct {
	template *template.Template
}

type Options struct {
	TemplateDirectory string
}

func NewHandler(options Options) (*Handler, error) {

	if len(options.TemplateDirectory) == 0 {
		options.TemplateDirectory = "pkg/webapp/templates/"
	}

	e, err := os.Executable()
	if err != nil {
		return nil, err
	}
	fmt.Println(e)

	t, err := template.ParseGlob(options.TemplateDirectory + "/*.gohtml")
	if err != nil {
		return nil, err
	}
	h := &Handler{
		template: t,
	}
	return h, nil
}
