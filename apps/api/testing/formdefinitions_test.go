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
    name        string
    description string
    formDef     *v1.FormDefinition
    customize   func(fd *v1.FormDefinition)
    assert      func(t *testing.T, fd *v1.FormDefinition, err error)
  }

  // defaultAssertErr will verify that an error response has the proper format
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
      name:        "Valid",
      description: "The form definition is valid and should not throw any error",
      formDef:     aValidFormDefinition("group", "Test1"),
      assert: func(t *testing.T, fd *v1.FormDefinition, err error) {
        assert.NoError(t, err)
      },
    }, {
      name:        "MissingSpecGroup",
      description: "The form definition is missing the .Spec.Group property, which is required",
      formDef:     aValidFormDefinition("group", "MissingSpecGroup"),
      customize:   func(fd *v1.FormDefinition) { fd.Spec.Group = "" },
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
      name:        "MissingPlural",
      description: "The form definition is missing the .Spec.Names.Plural property, which is required",
      formDef:     aValidFormDefinition("group", "MissingPlural"),
      customize:   func(fd *v1.FormDefinition) { fd.Spec.Names.Plural = "" },
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
      name:        "MissingSingular",
      description: "The form definition is missing the .Spec.Names.Singular property, which is required",
      formDef:     aValidFormDefinition("group", "MissingSingular"),
      customize:   func(fd *v1.FormDefinition) { fd.Spec.Names.Singular = "" },
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
      name:        "EmptyVersions",
      description: "The form definition does not specify any version. It must at least define one",
      formDef:     aValidFormDefinition("group", "EmptyVersions"),
      customize:   func(fd *v1.FormDefinition) { fd.Spec.Versions = []v1.FormDefinitionVersion{} },
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
      name:        "MissingVersionName",
      description: "The form definition does not specify the version name, which is required",
      formDef:     aValidFormDefinition("group", "MissingVersionName"),
      customize:   func(fd *v1.FormDefinition) { fd.Spec.Versions[0].Name = "" },
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

  // run each test case
  for _, testCase := range testCases {

    // use a copy of the testCase as the parallel tests will overwrite this
    tc := testCase

    s.T().Run(tc.name, func(t *testing.T) {
      t.Parallel() // running these tests in parallel to simulate concurrent requests

      fd := tc.formDef
      if tc.customize != nil {
        tc.customize(fd)
      }
      out, err := s.nrcClient.FormDefinitions().Create(s.ctx, fd)
      if tc.assert != nil {
        tc.assert(t, out, err)
      }
      if err == nil {
        // deleting the entity if it was created (cleanup)
        assert.NoError(t, s.nrcClient.FormDefinitions().Delete(s.ctx, out.GetUID(), metav1.DeleteOptions{}))
      }
    })
  }

}

// aValidFormDefinition returns a valid form definition used for testing.
// it will automatically populate the names based on the provided group/kind
func aValidFormDefinition(group, kind string) *v1.FormDefinition {
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
