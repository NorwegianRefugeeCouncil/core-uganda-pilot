package form

import (
	"github.com/nrc-no/core/pkg/i18n"
)

//Section describes a sub-section of a form. It is simply a list of Control Names.
//NB this is a quick and dirty solution because deadline; you should probably FIXME
type Section struct {
	Title        i18n.Strings `json:"title" bson:"title"`
	ControlNames []string     `json:"controlNames" bson:"controlNames"`
}
