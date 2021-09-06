package cms

import (
	"encoding/json"
	"github.com/nrc-no/core/pkg/form"
	"net/url"
	"time"
)

// Case describes a case form including relevant metadata
type Case struct {
	ID string `json:"id" bson:"id"`
	// CaseTypeID is required
	CaseTypeID string `json:"caseTypeId" bson:"caseTypeId"`
	// PartyID is required
	PartyID string `json:"partyId" bson:"partyId"`
	// TeamID is required
	TeamID string `json:"teamId" bson:"teamId"`
	// CreatorID is required
	CreatorID string `json:"creatorId" bson:"creatorId"`
	// ParentID refers to a "parent" case, if this case is a referral.
	// ParentID is optional
	ParentID   string `json:"parentId" bson:"parentId"`
	IntakeCase bool   `json:"intakeCase" bson:"intakeCase"`
	// Form is required
	Form form.Form `json:"form" bson:"form"`
	// FormData is optional
	FormData url.Values `json:"formData"`
	Done     bool       `json:"done" bson:"done"`
}

type CaseList struct {
	Items []*Case `json:"items" bson:"items"`
}

// CaseType contains the information needed to construct a case form as well as Team and PartyType IDs associated with the form.
type CaseType struct {
	ID             string    `json:"id" bson:"id"`
	Name           string    `json:"name" bson:"name"`
	PartyTypeID    string    `json:"partyTypeId" bson:"partyTypeId"`
	TeamID         string    `json:"teamId" bson:"teamId"`
	Form           form.Form `json:"form" bson:"form"`
	IntakeCaseType bool      `json:"intakeCaseType" bson:"intakeCaseType"`
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
