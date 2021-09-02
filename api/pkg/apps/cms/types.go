package cms

import (
	"encoding/json"
	"github.com/nrc-no/core/pkg/form"
	"net/url"
	"time"
)

// Case describes a case form including relevant metadata
type Case struct {
	ID         string        `json:"id" bson:"id"`
	CaseTypeID string        `json:"caseTypeId" bson:"caseTypeId"`
	PartyID    string        `json:"partyId" bson:"partyId"`
	Done       bool          `json:"done" bson:"done"`
	ParentID   string        `json:"parentId" bson:"parentId"`
	TeamID     string        `json:"teamId" bson:"teamId"`
	CreatorID  string        `json:"creatorId" bson:"creatorId"`
	Template   *CaseTemplate `json:"template" bson:"template"`
	IntakeCase bool          `json:"intakeCase" bson:"intakeCase"`
}

type CaseList struct {
	Items []*Case `json:"items" bson:"items"`
}

// CaseType contains the information needed to construct a case form as well as Team and PartyType IDs associated with the form.
type CaseType struct {
	ID             string        `json:"id" bson:"id"`
	Name           string        `json:"name" bson:"name"`
	PartyTypeID    string        `json:"partyTypeId" bson:"partyTypeId"`
	TeamID         string        `json:"teamId" bson:"teamId"`
	Template       *CaseTemplate `json:"template" bson:"template"`
	IntakeCaseType bool          `json:"intakeCaseType" bson:"intakeCaseType"`
}

type CaseTypeList struct {
	Items []*CaseType `json:"items" bson:"items"`
}

func (c *CaseType) String() string {
	return c.Name
}

func (c *CaseType) Pretty() string {
	b, err := json.MarshalIndent(c.Form, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (l *CaseTypeList) FindByID(id string) *CaseType {
	for _, caseType := range l.Items {
		if caseType.ID == id {
			return caseType
		}
	}
	return nil
}

type Comment struct {
	ID        string    `json:"id" bson:"id"`
	CaseID    string    `json:"caseId" bson:"caseId"`
	AuthorID  string    `json:"authorId" bson:"authorId"`
	Body      string    `json:"body" bson:"body"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type CommentList struct {
	Items []*Comment `json:"items"`
}

// Case templates
// https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-githubs-form-schema

// CaseTemplate contains a list of form elements used to construct a case form
type CaseTemplate struct {
	// FormElements is an ordered list of the elements found in the form
	FormElements []form.FormElement `json:"formElements" bson:"formElements"`
}

func (c *CaseTemplate) MarkAsReadonly() {
	var elems []form.FormElement
	for _, element := range c.FormElements {
		element.Readonly = true
		elems = append(elems, element)
	}
	c.FormElements = elems
}
