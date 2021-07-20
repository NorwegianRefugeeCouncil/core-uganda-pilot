package webapp

import (
	"fmt"
	"github.com/nrc-no/core/pkg/auth"
	"github.com/nrc-no/core/pkg/sessionmanager"
	"github.com/nrc-no/core/pkg/validation"
	"html"
	"html/template"
	"io"
	"net/http"
	"path"
)

// RenderInterface defines the methods available in the template
type RenderInterface interface {
	IsLoggedIn() bool
	Profile() (*Claims, error)
	Notifications() ([]*sessionmanager.Notification, error)
	ValidationClass(errorList validation.ErrorList) string
	ValidationFeedback(errorList validation.ErrorList, id string) template.HTML
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

func (r *RendererFactory) Notifications() ([]*sessionmanager.Notification, error) {
	return []*sessionmanager.Notification{}, nil
}

func (r *RendererFactory) ValidationClass(errorList validation.ErrorList) string {
	return ""
}

func (r *RendererFactory) ValidationFeedback(errorList validation.ErrorList, id string) template.HTML {
	return ""
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
func (r *RendererFactory) New(req *http.Request) *Renderer {
	renderer := &Renderer{
		req:            req,
		sessionManager: r.sessionManager,
	}
	renderer.template = WithRenderInterface(r.template, renderer)
	return renderer
}

// Renderer is the actual struct that will render templates
type Renderer struct {
	template       *template.Template
	req            *http.Request
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

func (r *Renderer) Notifications() ([]*sessionmanager.Notification, error) {
	return r.sessionManager.ConsumeNotifications(r.req)
}

func (r *Renderer) ValidationClass(errorList validation.ErrorList) string {
	if len(errorList) > 0 {
		return "is-invalid"
	}
	return "is-valid"
}

func (r *Renderer) ValidationFeedback(errorList validation.ErrorList, id string) template.HTML {
	if len(errorList) == 0 {
		return `<div class="valid-feedback">Looks good!</div>`
	}
	s := fmt.Sprintf(`<div id="%sFeedback" class="invalid-feedback">`, id)
	for i, e := range errorList {
		if i > 0 {
			s += `<br>`
		}
		s += html.EscapeString(e.Detail)
	}
	s += `</div>`
	return template.HTML(s)
}

// WithRenderInterface adds the RenderInterface methods to the template
func WithRenderInterface(t *template.Template, intf RenderInterface) *template.Template {
	return t.Funcs(map[string]interface{}{
		"IsLoggedIn":         intf.IsLoggedIn,
		"Profile":            intf.Profile,
		"Notifications":      intf.Notifications,
		"ValidationClass":    intf.ValidationClass,
		"ValidationFeedback": intf.ValidationFeedback,
	})
}
