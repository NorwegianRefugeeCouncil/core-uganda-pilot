package webapp

import (
	"github.com/nrc-no/core-kafka/pkg/auth"
	"html/template"
	"io"
	"net/http"
)

// RenderInterface defines the methods available in the template
type RenderInterface interface {
	IsLoggedIn() bool
}

// RendererFactory is a factory to create Renderer
type RendererFactory struct {
	template *template.Template
}

// RendererFactory must implement RenderInterface so that the methods
// are available at compile time
var _ RenderInterface = &RendererFactory{}

// IsLoggedIn is a mock method that allows compilation of go templates
// Actual implementation is in Renderer.IsLoggedIn
func (r *RendererFactory) IsLoggedIn() bool {
	return false
}

// NewRendererFactory creates a new instance of the RendererFactory
func NewRendererFactory(templateDirectory string) (*RendererFactory, error) {
	f := &RendererFactory{}
	t := template.New("")
	t = WithRenderInterface(t, f)
	var err error
	t, err = t.ParseGlob(templateDirectory + "/*.gohtml")
	if err != nil {
		return nil, err
	}
	f.template = t
	return f, nil
}

// New creates a new Renderer
func (r *RendererFactory) New(req *http.Request) *Renderer {
	renderer := &Renderer{
		req: req,
	}
	renderer.template = WithRenderInterface(r.template, renderer)
	return renderer
}

// Renderer is the actual struct that will render templates
type Renderer struct {
	template *template.Template
	req      *http.Request
}

// Renderer must implement RenderInterface so that the methods are available
// in the templates
var _ RenderInterface = &Renderer{}

func (r *Renderer) ExecuteTemplate(w io.Writer, name string, args interface{}) error {
	return r.template.ExecuteTemplate(w, name, args)
}

// IsLoggedIn returns whether the request is made by an authenticated user
func (r *Renderer) IsLoggedIn() bool {
	return auth.IsAuthenticatedRequest(r.req)
}

// WithRenderInterface adds the RenderInterface methods to the template
func WithRenderInterface(t *template.Template, intf RenderInterface) *template.Template {
	return t.Funcs(map[string]interface{}{
		"IsLoggedIn": intf.IsLoggedIn,
	})
}
