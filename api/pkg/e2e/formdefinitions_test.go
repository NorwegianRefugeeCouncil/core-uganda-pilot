package e2e

import (
	"context"
	corev1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	"github.com/stretchr/testify/assert"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	var found = false
	for _, cause := range status.Details.Causes {
		if cause.Type == causeType {
			found = true
			assert.Equal(t, field, cause.Field)
			assert.Equal(t, message, cause.Message)
			break
		}
	}

	if !found {
		t.Fatalf("could not find cause %s in response. Actual response had %v causes", causeType, status.Details.Causes)
	}

}

func (s *Suite) TestCreateFormDefinition() {
	t := s.T()

	ctx := context.TODO()

	// Predefined assertion that ensures that no error is returned
	var assertNoError = func(t *testing.T, in, out *corev1.FormDefinition, err error) {
		assert.NoError(t, err)
	}

	// Predefined assertion that asserts the structure of an error: the type, the involved field path and the message
	var assertStatusCause = func(cause metav1.CauseType, field, message string) func(t *testing.T, in, out *corev1.FormDefinition, err error) {
		return func(t *testing.T, in, out *corev1.FormDefinition, err error) {
			assert.Error(t, err)
			assert.True(t, apierrors.IsInvalid(err))
			assertHasStatusCause(t, err, cause, field, message)
		}
	}

	var testCases = []struct {
		name      string
		customize func(f *corev1.FormDefinition)
		assert    func(t *testing.T, in, out *corev1.FormDefinition, err error)
	}{
		{
			name:   "valid",
			assert: assertNoError,
		}, {
			name: "missing group",
			customize: func(f *corev1.FormDefinition) {
				f.Spec.Group = ""
			},
			assert: assertStatusCause(metav1.CauseTypeFieldValueRequired, "spec.group", "Required value: group is required"),
		}, {
			name: "missing plural name",
			customize: func(f *corev1.FormDefinition) {
				f.Spec.Names.Plural = ""
			},
			assert: assertStatusCause(metav1.CauseTypeFieldValueRequired, "spec.names.plural", "Required value: plural is required"),
		}, {
			name: "missing singular name",
			customize: func(f *corev1.FormDefinition) {
				f.Spec.Names.Singular = ""
			},
			assert: assertStatusCause(metav1.CauseTypeFieldValueRequired, "spec.names.singular", "Required value: singular is required"),
		}, {
			name: "missing kind name",
			customize: func(f *corev1.FormDefinition) {
				f.Spec.Names.Kind = ""
			},
			assert: assertStatusCause(metav1.CauseTypeFieldValueRequired, "spec.names.kind", "Required value: kind is required"),
		}, {
			name: "empty versions",
			customize: func(f *corev1.FormDefinition) {
				f.Spec.Versions = []corev1.FormDefinitionVersion{}
			},
			assert: assertStatusCause(metav1.CauseTypeFieldValueRequired, "spec.versions", "Required value: versions cannot be empty"),
		}, {
			name: "empty version name",
			customize: func(f *corev1.FormDefinition) {
				f.Spec.Versions[0].Name = ""
			},
			assert: assertStatusCause(metav1.CauseTypeFieldValueRequired, "spec.versions[0].name", "Required value: version name cannot be empty"),
		}, {
			name: "duplicate version name",
			customize: func(f *corev1.FormDefinition) {
				versionCopy := f.Spec.Versions[0]
				f.Spec.Versions = append(f.Spec.Versions, versionCopy)
			},
			assert: assertStatusCause(metav1.CauseTypeFieldValueDuplicate, "spec.versions[1].name", "Duplicate value: \"v1\""),
		}, {
			name: "missing key",
			customize: func(f *corev1.FormDefinition) {
				f.Spec.Versions[0].Schema.FormSchema.Root.Children[0].Key = ""
			},
			assert: assertStatusCause(metav1.CauseTypeFieldValueRequired, "spec.versions[0].schema.formSchema.root.children[0].key", "Required value: key is required"),
		}, {
			name: "root has key defined",
			customize: func(f *corev1.FormDefinition) {
				f.Spec.Versions[0].Schema.FormSchema.Root.Key = "shouldNotBeThere"
			},
			assert: assertStatusCause(metav1.CauseTypeFieldValueNotSupported, "spec.versions[0].schema.formSchema.root.key", "Unsupported value: \"shouldNotBeThere\": supported values: \"\""),
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fd := aValidFormDefinition()
			if tc.customize != nil {
				tc.customize(fd)
			}

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
			Group: "test",
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
