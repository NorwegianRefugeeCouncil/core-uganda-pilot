package e2e

import (
	"context"
	corev1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	"github.com/stretchr/testify/assert"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"testing"
)

func (s *Suite) TestListFormDefinitions() {
	t := s.T()
	ctx := context.TODO()
	_, err := s.client.FormDefinitions().List(ctx, metav1.ListOptions{})
	assert.NoError(t, err)
}

func assertHasStatusCause(t *testing.T, err error, causeType metav1.CauseType, field, message string) {
	apierr, ok := err.(apierrors.APIStatus)
	if !assert.True(t, ok) {
		return
	}

	status := apierr.Status()
	if !assert.NotNil(t, status.Details) {
		return
	}

	for _, cause := range status.Details.Causes {
		if cause.Type == causeType && cause.Field == field && cause.Message == message {
			return
		}
	}

	t.Errorf("could not find cause %s in response. Actual response had %v causes", causeType, status.Details.Causes)

}

func (s *Suite) TestCreateFormDefinition() {
	t := s.T()

	ctx := context.TODO()

	type assertionFunc func(t *testing.T, in, out *corev1.FormDefinition, err error)

	// Predefined assertion that ensures that no error is returned
	var assertNoError = func(t *testing.T, in, out *corev1.FormDefinition, err error) {
		assert.NoError(t, err)
	}

	// Predefined assertion that asserts the structure of an error: the type, the involved field path and the message
	var assertStatusCause = func(cause metav1.CauseType, field, message string) assertionFunc {
		return func(t *testing.T, in, out *corev1.FormDefinition, err error) {
			assert.Error(t, err)
			assert.True(t, apierrors.IsInvalid(err))
			assertHasStatusCause(t, err, cause, field, message)
		}
	}

	// Predefined assertions to combine multiple assertions together
	var assertMultiple = func(assertions ...assertionFunc) assertionFunc {
		return func(t *testing.T, in, out *corev1.FormDefinition, err error) {
			for _, assertion := range assertions {
				assertion(t, in, out, err)
			}
		}
	}

	var int64ptr = func(i int64) *int64 {
		return &i
	}

	var testCases = []struct {
		name      string
		customize func(f *corev1.FormDefinition)
		assert    func(t *testing.T, in, out *corev1.FormDefinition, err error)
	}{
		{ // a valid form
			name:   "valid",
			assert: assertNoError,
		}, { // the group is missing
			name: "missing group",
			customize: func(f *corev1.FormDefinition) {
				f.Spec.Group = ""
			},
			assert: assertStatusCause(metav1.CauseTypeFieldValueRequired, "spec.group", "Required value: group is required"),
		}, { // the plural name is missing
			name: "missing plural name",
			customize: func(f *corev1.FormDefinition) {
				f.Spec.Names.Plural = ""
			},
			assert: assertStatusCause(metav1.CauseTypeFieldValueRequired, "spec.names.plural", "Required value: plural is required"),
		}, { // the singular name is missing
			name: "missing singular name",
			customize: func(f *corev1.FormDefinition) {
				f.Spec.Names.Singular = ""
			},
			assert: assertStatusCause(metav1.CauseTypeFieldValueRequired, "spec.names.singular", "Required value: singular is required"),
		}, { // the kind name is missing
			name: "missing kind name",
			customize: func(f *corev1.FormDefinition) {
				f.Spec.Names.Kind = ""
			},
			assert: assertStatusCause(metav1.CauseTypeFieldValueRequired, "spec.names.kind", "Required value: kind is required"),
		}, { // empty versions array. should not allow to post a form definition with no version at all
			name: "empty versions",
			customize: func(f *corev1.FormDefinition) {
				f.Spec.Versions = []corev1.FormDefinitionVersion{}
			},
			assert: assertStatusCause(metav1.CauseTypeFieldValueRequired, "spec.versions", "Required value: versions cannot be empty"),
		}, { // empty version name
			name: "empty version name",
			customize: func(f *corev1.FormDefinition) {
				f.Spec.Versions[0].Name = ""
			},
			assert: assertStatusCause(metav1.CauseTypeFieldValueRequired, "spec.versions[0].name", "Required value: version name cannot be empty"),
		}, { // two versions with the same name
			name: "duplicate version name",
			customize: func(f *corev1.FormDefinition) {
				versionCopy := f.Spec.Versions[0]
				f.Spec.Versions = append(f.Spec.Versions, versionCopy)
			},
			assert: assertStatusCause(metav1.CauseTypeFieldValueDuplicate, "spec.versions[1].name", "Duplicate value: \"v1\""),
		}, { // missing key on a form element that requires a key
			name: "missing key",
			customize: func(f *corev1.FormDefinition) {
				f.Spec.Versions[0].Schema.FormSchema.Root.Children[0].Key = ""
			},
			assert: assertStatusCause(metav1.CauseTypeFieldValueRequired, "spec.versions[0].schema.formSchema.root.children[0].key", "Required value: key is required"),
		}, { // root key cannot be defined
			name: "root has key defined",
			customize: func(f *corev1.FormDefinition) {
				f.Spec.Versions[0].Schema.FormSchema.Root.Key = "shouldNotBeThere"
			},
			assert: assertStatusCause(metav1.CauseTypeFieldValueNotSupported, "spec.versions[0].schema.formSchema.root.key", "Unsupported value: \"shouldNotBeThere\": supported values: \"\""),
		}, { // text inputs should allow minLength/maxLength validation property
			name: "textinput should allow minLength/maxLength",
			customize: func(f *corev1.FormDefinition) {
				children := f.Spec.Versions[0].Schema.FormSchema.Root.Children
				f.Spec.Versions[0].Schema.FormSchema.Root.Children = append(children, corev1.FormElementDefinition{
					Type:      corev1.ShortTextType,
					Key:       "somekey",
					MinLength: 3,
					MaxLength: int64ptr(4),
				})
			},
			assert: assertNoError,
		}, { // text inputs should not allow minLength > maxLength
			name: "textinput should not allow minLength > maxLength",
			customize: func(f *corev1.FormDefinition) {
				children := f.Spec.Versions[0].Schema.FormSchema.Root.Children
				f.Spec.Versions[0].Schema.FormSchema.Root.Children = append(children, corev1.FormElementDefinition{
					Type:      corev1.ShortTextType,
					Key:       "somekey",
					MinLength: 10,
					MaxLength: int64ptr(4),
				})
			},
			assert: assertMultiple(
				assertStatusCause(
					metav1.CauseTypeFieldValueInvalid,
					"spec.versions[0].schema.formSchema.root.children[2].minLength",
					"Invalid value: 10: maximum length cannot be smaller than minimum length"),
				assertStatusCause(
					metav1.CauseTypeFieldValueInvalid,
					"spec.versions[0].schema.formSchema.root.children[2].maxLength",
					"Invalid value: 4: maximum length cannot be smaller than minimum length"),
			),
		}, { // textInput cannot accept min/max parameters
			name: "textinput should not allow min",
			customize: func(f *corev1.FormDefinition) {
				children := f.Spec.Versions[0].Schema.FormSchema.Root.Children
				f.Spec.Versions[0].Schema.FormSchema.Root.Children = append(children, corev1.FormElementDefinition{
					Type: corev1.ShortTextType,
					Key:  "somekey",
					Min:  "2",
				})
			},
			assert: assertStatusCause(
				metav1.CauseTypeFieldValueNotSupported,
				"spec.versions[0].schema.formSchema.root.children[2].min",
				"Unsupported value: \"2\": supported values: \"\""),
		}, { // textInput cannot accept min/max parameters
			name: "textinput should not allow max",
			customize: func(f *corev1.FormDefinition) {
				children := f.Spec.Versions[0].Schema.FormSchema.Root.Children
				f.Spec.Versions[0].Schema.FormSchema.Root.Children = append(children, corev1.FormElementDefinition{
					Type: corev1.ShortTextType,
					Key:  "somekey",
					Max:  "2",
				})
			},
			assert: assertStatusCause(
				metav1.CauseTypeFieldValueNotSupported,
				"spec.versions[0].schema.formSchema.root.children[2].max",
				"Unsupported value: \"2\": supported values: \"\""),
		}, { // integer inputs should allow min/max validation parameters
			name: "integer input should allow min/max",
			customize: func(f *corev1.FormDefinition) {
				children := f.Spec.Versions[0].Schema.FormSchema.Root.Children
				f.Spec.Versions[0].Schema.FormSchema.Root.Children = append(children, corev1.FormElementDefinition{
					Type: corev1.IntegerType,
					Key:  "somekey",
					Min:  "0",
					Max:  "2",
				})
			},
			assert: assertNoError,
		}, { // integer inputs should not allow min > max
			name: "integer input should not allow min > max",
			customize: func(f *corev1.FormDefinition) {
				children := f.Spec.Versions[0].Schema.FormSchema.Root.Children
				f.Spec.Versions[0].Schema.FormSchema.Root.Children = append(children, corev1.FormElementDefinition{
					Type: corev1.IntegerType,
					Key:  "somekey",
					Min:  "2",
					Max:  "0",
				})
			},
			assert: assertMultiple(
				assertStatusCause(
					metav1.CauseTypeFieldValueInvalid,
					"spec.versions[0].schema.formSchema.root.children[2].min",
					"Invalid value: \"2\": minimum cannot be greater than maximum"),
				assertStatusCause(
					metav1.CauseTypeFieldValueInvalid,
					"spec.versions[0].schema.formSchema.root.children[2].max",
					"Invalid value: \"0\": minimum cannot be greater than maximum"),
			),
		}, { // integer inputs cannot accept invalid numerical values for min/max
			name: "integer input should not allow invalid min/max",
			customize: func(f *corev1.FormDefinition) {
				children := f.Spec.Versions[0].Schema.FormSchema.Root.Children
				f.Spec.Versions[0].Schema.FormSchema.Root.Children = append(children, corev1.FormElementDefinition{
					Type: corev1.IntegerType,
					Key:  "somekey",
					Min:  "abc",
					Max:  "def",
				})
			},
			assert: assertMultiple(
				assertStatusCause(
					metav1.CauseTypeFieldValueInvalid,
					"spec.versions[0].schema.formSchema.root.children[2].min",
					"Invalid value: \"abc\": invalid number"),
				assertStatusCause(
					metav1.CauseTypeFieldValueInvalid,
					"spec.versions[0].schema.formSchema.root.children[2].max",
					"Invalid value: \"def\": invalid number"),
			),
		},
	}

	for i, testCase := range testCases {
		tc := testCase
		var idx = i
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fd := aValidFormDefinition()
			if tc.customize != nil {
				tc.customize(fd)
			}
			fd.Name = "test-" + strconv.Itoa(idx)

			err := s.client.FormDefinitions().Delete(ctx, fd.Name, metav1.DeleteOptions{})
			if err != nil && !apierrors.IsNotFound(err) {
				assert.NoError(t, err)
				return
			}

			out, err := s.client.FormDefinitions().Create(ctx, fd, metav1.CreateOptions{})
			tc.assert(t, fd, out, err)

		})
	}
}

func aValidFormDefinition() *corev1.FormDefinition {
	return &corev1.FormDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
		Spec: corev1.FormDefinitionSpec{
			Group: "test.com",
			Names: corev1.FormDefinitionNames{
				Plural:   "formtests",
				Singular: "formtest",
				Kind:     "FormTest",
			},
			Versions: []corev1.FormDefinitionVersion{
				{
					Name: "v1",
					Schema: corev1.FormDefinitionValidation{
						FormSchema: corev1.FormDefinitionSchema{
							Root: corev1.FormElementDefinition{
								Key:  "",
								Type: corev1.SectionType,
								Children: []corev1.FormElementDefinition{
									{
										Key: "firstName",
										Label: corev1.TranslatedStrings{
											{
												Locale: "en",
												Value:  "First Name",
											}, {
												Locale: "fr",
												Value:  "Prenom",
											},
										},
										Description: corev1.TranslatedStrings{
											{
												Locale: "en",
												Value:  "Enter the first name of the beneficiary",
											}, {
												Locale: "fr",
												Value:  "Entrez le prénom du bénéficiaire",
											},
										},
										Type:     corev1.ShortTextType,
										Required: true,
										Children: nil,
									}, {
										Key: "lastName",
										Label: corev1.TranslatedStrings{
											{
												Locale: "en",
												Value:  "Last Name",
											}, {
												Locale: "fr",
												Value:  "Nom de famille",
											},
										},
										Description: corev1.TranslatedStrings{
											{
												Locale: "en",
												Value:  "Enter the first name of the beneficiary",
											}, {
												Locale: "fr",
												Value:  "Entrez le nom de famille du bénéficiaire",
											},
										},
										Type:     corev1.ShortTextType,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
