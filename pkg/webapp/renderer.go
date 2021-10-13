package webapp

import (
	"fmt"
	"github.com/nrc-no/core/pkg/auth"
	"github.com/nrc-no/core/pkg/sessionmanager"
	"html/template"
	"io"
	"net/http"
	"path"
)

// RenderInterface defines the methods available in the template
type RenderInterface interface {
	IsLoggedIn() bool
	Profile() (*Claims, error)
}

// RendererFactory is a factory to create Renderer
type RendererFactory struct {
	template       *template.Template
	sessionManager sessionmanager.Store
}

// RendererFactory must implement RenderInterface so that the methods
// are available at compile time
var _ RenderInterface = &RendererFactory{}

// IsLoggedIn is a mock method that allows compilation of go templates
// Actual implementation is in Renderer.IsLoggedIn
func (r *RendererFactory) IsLoggedIn() bool {
	return false
}

func (r *RendererFactory) Profile() (*Claims, error) {
	return nil, nil
}

// NewRendererFactory creates a new instance of the RendererFactory
func NewRendererFactory(templateDirectory string, sessionManager sessionmanager.Store) (*RendererFactory, error) {
	f := &RendererFactory{
		sessionManager: sessionManager,
	}
	t := template.New("")
	t = WithRenderInterface(t, f)

	var err error

	t, err = t.ParseGlob(path.Join(templateDirectory, "*.gohtml"))
	if err != nil {
		return nil, err
	}

	f.template = t

	return f, nil
}

// New creates a new Renderer
func (r *RendererFactory) New(req *http.Request, w http.ResponseWriter) *Renderer {
	renderer := &Renderer{
		req:            req,
		w:              w,
		sessionManager: r.sessionManager,
	}
	renderer.template = WithRenderInterface(r.template, renderer)

	return renderer
}

// Renderer is the actual struct that will render templates
type Renderer struct {
	template       *template.Template
	req            *http.Request
	w              http.ResponseWriter
	sessionManager sessionmanager.Store
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

func (r *Renderer) Profile() (*Claims, error) {
	session, err := r.sessionManager.Get(r.req)
	if err != nil {
		return nil, err
	}

	profileIntf, ok := session.Values["profile"]
	if !ok {
		return nil, fmt.Errorf("profile not found")
	}

	profile, ok := profileIntf.(*Claims)
	if !ok {
		return nil, fmt.Errorf("")
	}

	return profile, nil
}

// WithRenderInterface adds the RenderInterface methods to the template
func WithRenderInterface(t *template.Template, intf RenderInterface) *template.Template {
	return t.Funcs(map[string]interface{}{
		"IsLoggedIn": intf.IsLoggedIn,
		"Profile":    intf.Profile,
	})
}
