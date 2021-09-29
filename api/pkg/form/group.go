package form

//Group describes a sub-section of a form. It is simply a list of Control Names.
//NB this is a quick and dirty solution because deadline; you should probably FIXME
type Group struct {
	Controls []string `json:"controls" bson:"controls"`
}
