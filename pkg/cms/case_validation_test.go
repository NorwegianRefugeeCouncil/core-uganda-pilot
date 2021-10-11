package cms_test

import (
	"github.com/nrc-no/core/internal/validation"
	. "github.com/nrc-no/core/pkg/cms"
	"github.com/stretchr/testify/assert"
	"testing"
)

var names = []string{"text", "email", "phone", "url", "date", "textarea", "dropdown", "checkbox", "radio", "taxonomy"}
var caseType = mockCaseType()

func TestValidateCase(t *testing.T) {

	tcs := []struct {
		name   string
		kase   Case
		assert func(t *testing.T, errList validation.ErrorList)
	}{{
		name: "missing required fields",
		kase: Case{
			CaseTypeID: caseType.ID,
			Form:       caseType.Form,
			FormData: map[string][]string{
				"text":     nil,
				"email":    nil,
				"phone":    nil,
				"url":      nil,
				"date":     nil,
				"textarea": nil,
				"dropdown": nil,
				"checkbox": nil,
				"radio":    nil,
				"taxonomy": nil,
				"file":     nil,
			},
		},
		assert: func(t *testing.T, errList validation.ErrorList) {
			assert.NotEmpty(t, errList)
			errLists := []*validation.ErrorList{}
			for _, name := range names {
				errLists = append(errLists, errList.Find(name))
			}
			for _, list := range errLists {
				assert.NotNil(t, list)
				assert.NotEmpty(t, list)
				assert.Len(t, *list, 1)
				l := *list
				err := l[0]
				assert.Equal(t, err.Type, validation.ErrorTypeRequired)
			}
		},
	},
		{
			name: "valid fields",
			kase: Case{
				ID:         newUUID(),
				CaseTypeID: caseType.ID,
				PartyID:    newUUID(),
				TeamID:     newUUID(),
				CreatorID:  newUUID(),
				ParentID:   newUUID(),
				Form:       caseType.Form,
				FormData: map[string][]string{
					"text":     {"test"},
					"email":    {"test@email.com"},
					"phone":    {"0555-555555"},
					"url":      {"https://www.example.com"},
					"date":     {"1967-03-23"},
					"textarea": {"test"},
					"dropdown": {"test"},
					"checkbox": {"test"},
					"radio":    {"test"},
					"taxonomy": {"test"},
					"file":     {"test"},
				},
			},
			assert: func(t *testing.T, errList validation.ErrorList) {
				assert.Empty(t, errList)
			},
		},
	}

	for _, tc := range tcs {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			errList := ValidateCase(&testCase.kase, validation.NewPath(""))
			testCase.assert(t, errList)
		})
	}

}
