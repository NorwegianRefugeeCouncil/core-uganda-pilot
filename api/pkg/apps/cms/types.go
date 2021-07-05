package cms

import (
	"encoding/json"
	"net/url"
	"time"
)

type Case struct {
	ID          string `json:"id" bson:"id"`
	CaseTypeID  string `json:"caseTypeId" bson:"caseTypeId"`
	PartyID     string `json:"partyId" bson:"partyId"`
	Description string `json:"description" bson:"description"`
	Done        bool   `json:"done" bson:"done"`
	ParentID    string `json:"parentId" bson:"parentId"`
	TeamID      string `json:"teamId" bson:"teamId"`
	CreatorID   string `json:"creatorId" bson:"creatorId"`
}

type CaseList struct {
	Items []*Case `json:"items" bson:"items"`
}

type CaseType struct {
	ID          string       `json:"id" bson:"id"`
	Name        string       `json:"name" bson:"name"`
	PartyTypeID string       `json:"partyTypeId" bson:"partyTypeId"`
	TeamID      string       `json:"teamId" bson:"teamId"`
	Template    CaseTemplate `json:"template" bson:"template"`
}

type CaseTypeList struct {
	Items []*CaseType `json:"items" bson:"items"`
}

func (c *CaseType) String() string {
	return c.Name
}

func (c *CaseType) UnmarshalFormData(values url.Values) error {
	c.Name = values.Get("name")
	c.PartyTypeID = values.Get("partyTypeId")
	c.TeamID = values.Get("teamId")
	templateString := values.Get("template")
	var formElements []CaseTemplateFormElement
	if err := json.Unmarshal([]byte(templateString), &formElements); err != nil {
		return err
	}
	c.Template = CaseTemplate{FormElements: formElements}
	return nil
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

type CaseTemplate struct {
	// FormElements is an ordered list of the elements found in the form
	FormElements []CaseTemplateFormElement `json:"formElements" bson:"formElements"`
}

type CaseTemplateFormElement struct {
	Type       string                            `json:"type" bson:"type"`
	Attributes CaseTemplateFormElementAttribute  `json:"attributes" bson:"attributes"`
	Validation CaseTemplateFormElementValidation `json:"validation" bson:"validation"`
}

type CaseTemplateFormElementAttribute struct {
	Label           string                       `json:"label" bson:"label"`
	ID              string                       `json:"id" bson:"id"`
	Description     string                       `json:"description" bson:"description"`
	Placeholder     string                       `json:"placeholder" bson:"placeholder"`
	Value           string                       `json:"value" bson:"value"`
	Multiple        bool                         `json:"multiple" bson:"multiple"`
	Options         []string                     `json:"options" bson:"options"`
	CheckboxOptions []CaseTemplateCheckboxOption `json:"checkboxOptions" bson:"checkboxOptions"`
}

type CaseTemplateFormElementValidation struct {
	Required bool `json:"required" bson:"required"`
}

type CaseTemplateCheckboxOption struct {
	Label    string `json:"label" bson:"label"`
	Required bool   `json:"required" bson:"required"`
}
