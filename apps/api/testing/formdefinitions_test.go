package testing

import (
	"github.com/nrc-no/core/apps/api/pkg/apis/core/v1"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func (s *MainTestSuite) TestFormDefinitionPostValidationErrors() {

	type testCase struct {
		name           string
		formDefinition *v1.FormDefinition
		customizeFn    func(fd *v1.FormDefinition)
		assert         func(t *testing.T, fd *v1.FormDefinition, err error)
	}

	var defaultAssertErr = func(t *testing.T, err error, causeCount int, out *metav1.Status) bool {
		var status = &metav1.Status{}
		if !AssertIsErrStatus(t, err, status) {
			return false
		}
		AssertStatusFailure(t, status)
		AssertStatusCode(t, status, http.StatusUnprocessableEntity)
		AssertFailureReason(t, status, metav1.StatusReasonInvalid)
		AssertStatusMessage(t, status, "could not process entity")
		AssertStatusDetailsKind(t, status, "FormDefinition")
		AssertStatusDetailsGroup(t, status, "core")

		*out = *status
		if !AssertStatusCauseCount(t, status, causeCount) {
			return false
		}

		return true
	}

	var testCases = []testCase{
		{
			name:           "Valid",
			formDefinition: testFormDef("group", "Test1"),
			assert: func(t *testing.T, fd *v1.FormDefinition, err error) {
				assert.NoError(t, err)
			},
		}, {
			name:           "MissingSpecGroup",
			formDefinition: testFormDef("group", "MissingSpecGroup"),
			customizeFn:    func(fd *v1.FormDefinition) { fd.Spec.Group = "" },
			assert: func(t *testing.T, fd *v1.FormDefinition, err error) {
				var status = &metav1.Status{}
				if !defaultAssertErr(t, err, 1, status) {
					return
				}
				cause := status.Details.Causes[0]
				AssertCauseMessage(t, cause, "Required value: group is required")
				AssertCauseType(t, cause, metav1.CauseTypeFieldValueRequired)
				AssertCauseField(t, cause, "spec.group")
			},
		}, {
			name:           "MissingAPIVersion",
			formDefinition: testFormDef("group", "MissingAPIVersion"),
			customizeFn:    func(fd *v1.FormDefinition) { fd.APIVersion = "" },
			assert: func(t *testing.T, fd *v1.FormDefinition, err error) {
				var status = &metav1.Status{}
				if !defaultAssertErr(t, err, 1, status) {
					return
				}
				cause := status.Details.Causes[0]
				AssertCauseMessage(t, cause, "Required value: apiVersion is required")
				AssertCauseType(t, cause, metav1.CauseTypeFieldValueRequired)
				AssertCauseField(t, cause, "apiVersion")
			},
		}, {
			name:           "MissingKind",
			formDefinition: testFormDef("group", "kind"),
			customizeFn:    func(fd *v1.FormDefinition) { fd.Spec.Names.Kind = "" },
			assert: func(t *testing.T, fd *v1.FormDefinition, err error) {
				var status = &metav1.Status{}
				if !defaultAssertErr(t, err, 1, status) {
					return
				}
				cause := status.Details.Causes[0]
				AssertCauseMessage(t, cause, "Required value: kind name is required")
				AssertCauseType(t, cause, metav1.CauseTypeFieldValueRequired)
				AssertCauseField(t, cause, "spec.names.kind")
			},
		}, {
			name:           "MissingPlural",
			formDefinition: testFormDef("group", "MissingPlural"),
			customizeFn:    func(fd *v1.FormDefinition) { fd.Spec.Names.Plural = "" },
			assert: func(t *testing.T, fd *v1.FormDefinition, err error) {
				var status = &metav1.Status{}
				if !defaultAssertErr(t, err, 1, status) {
					return
				}
				cause := status.Details.Causes[0]
				AssertCauseMessage(t, cause, "Required value: plural name is required")
				AssertCauseType(t, cause, metav1.CauseTypeFieldValueRequired)
				AssertCauseField(t, cause, "spec.names.plural")
			},
		}, {
			name:           "MissingSingular",
			formDefinition: testFormDef("group", "MissingSingular"),
			customizeFn:    func(fd *v1.FormDefinition) { fd.Spec.Names.Singular = "" },
			assert: func(t *testing.T, fd *v1.FormDefinition, err error) {
				var status = &metav1.Status{}
				if !defaultAssertErr(t, err, 1, status) {
					return
				}
				cause := status.Details.Causes[0]
				AssertCauseMessage(t, cause, "Required value: singular name is required")
				AssertCauseType(t, cause, metav1.CauseTypeFieldValueRequired)
				AssertCauseField(t, cause, "spec.names.singular")
			},
		}, {
			name:           "EmptyVersions",
			formDefinition: testFormDef("group", "EmptyVersions"),
			customizeFn:    func(fd *v1.FormDefinition) { fd.Spec.Versions = []v1.FormDefinitionVersion{} },
			assert: func(t *testing.T, fd *v1.FormDefinition, err error) {
				var status = &metav1.Status{}
				if !defaultAssertErr(t, err, 1, status) {
					return
				}
				cause := status.Details.Causes[0]
				AssertCauseMessage(t, cause, "Required value: versions must not be empty")
				AssertCauseType(t, cause, metav1.CauseTypeFieldValueRequired)
				AssertCauseField(t, cause, "spec.versions")
			},
		}, {
			name:           "MissingVersionName",
			formDefinition: testFormDef("group", "MissingVersionName"),
			customizeFn:    func(fd *v1.FormDefinition) { fd.Spec.Versions[0].Name = "" },
			assert: func(t *testing.T, fd *v1.FormDefinition, err error) {
				var status = &metav1.Status{}
				if !defaultAssertErr(t, err, 1, status) {
					return
				}
				cause := status.Details.Causes[0]
				AssertCauseMessage(t, cause, "Required value: version name is required")
				AssertCauseType(t, cause, metav1.CauseTypeFieldValueRequired)
				AssertCauseField(t, cause, "spec.versions[0].name")
			},
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.name, func(t *testing.T) {
			fd := testCase.formDefinition
			if testCase.customizeFn != nil {
				testCase.customizeFn(fd)
			}
			out, err := s.nrcClient.FormDefinitions().Create(s.ctx, fd)
			if testCase.assert != nil {
				testCase.assert(t, out, err)
			}
		})
	}
}

func testFormDef(group, kind string) *v1.FormDefinition {
	return &v1.FormDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "core/v1",
			Kind:       "FormDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: v1.FormDefinitionSpec{
			Group: group,
			Names: v1.CustomResourceNames{
				Plural:   strings.ToLower(kind) + "s",
				Singular: strings.ToLower(kind),
				Kind:     kind,
			},
			Versions: []v1.FormDefinitionVersion{
				{Name: "v1", Schema: v1.FormSchema{
					FormSchema: v1.FormSchemaDefinition{
						Root: v1.FormElement{
							Key:         "root",
							ID:          uuid.NewV4().String(),
							Name:        []v1.TranslatedString{},
							Type:        "section",
							Description: []v1.TranslatedString{},
							Children: []v1.FormElement{
								{
									Key: "child",
									ID:  uuid.NewV4().String(),
									Name: []v1.TranslatedString{
										{Locale: "en", Value: "Name"},
									},
									Type: "text",
									Description: []v1.TranslatedString{
										{Locale: "en", Value: "Description"},
									},
									Children: []v1.FormElement{},
								},
							},
						},
					},
				}},
			},
		},
	}
}
