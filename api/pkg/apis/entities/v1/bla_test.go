package v1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func exprPtr(expr Expression) *Expression {
	return &expr
}

func TestValidation(t *testing.T) {

	testCases := []struct {
		name    string
		payload map[string]interface{}
		fields  []DataField
		assert  func(t *testing.T, form Form, submission FormSubmission, err error)
	}{
		{
			name: "required textual field: pass if payload property is present",
			payload: map[string]interface{}{
				"firstName": "John",
			},
			fields: []DataField{
				{
					Key:       "firstName",
					FieldType: ShortText,
					Required:  &TrueExpression,
				},
			},
			assert: func(t *testing.T, form Form, submission FormSubmission, err error) {
				assert.NoError(t, err)
			},
		}, {
			name: "required textual field: fail if payload property is missing",
			payload: map[string]interface{}{
				"firstName": "",
			},
			fields: []DataField{
				{
					Key:       "firstName",
					FieldType: ShortText,
					Required:  &TrueExpression,
				},
			},
			assert: func(t *testing.T, form Form, submission FormSubmission, err error) {
				assert.Error(t, err)
			},
		}, {
			name: "required textual field: fail if the payload property contains only spaces",
			payload: map[string]interface{}{
				"firstName": "  ",
			},
			fields: []DataField{
				{
					Key:       "firstName",
					FieldType: ShortText,
					Required:  &TrueExpression,
				},
			},
			assert: func(t *testing.T, form Form, submission FormSubmission, err error) {
				assert.Error(t, err)
			},
		}, {
			name:    "required textual field: fail if the payload property is missing",
			payload: map[string]interface{}{},
			fields: []DataField{
				{
					Key:       "firstName",
					FieldType: ShortText,
					Required:  &TrueExpression,
				},
			},
			assert: func(t *testing.T, form Form, submission FormSubmission, err error) {
				assert.Error(t, err)
			},
		}, {
			name: "required textual field: fail if the payload property is of wrong data type",
			payload: map[string]interface{}{
				"firstName": 123,
			},
			fields: []DataField{
				{
					Key:       "firstName",
					FieldType: ShortText,
					Required:  &TrueExpression,
				},
			},
			assert: func(t *testing.T, form Form, submission FormSubmission, err error) {
				assert.Error(t, err)
			},
		}, {
			name:    "required textual field with default value: should not throw if the payload property is missing",
			payload: map[string]interface{}{},
			fields: []DataField{
				{
					Key:       "firstName",
					FieldType: ShortText,
					Required:  &TrueExpression,
					Default:   exprPtr("'default'"),
				},
			},
			assert: func(t *testing.T, form Form, submission FormSubmission, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "default", submission.Payload["firstName"])
			},
		}, {
			name: "required numeric field: should not throw if the payload property is present",
			fields: []DataField{
				{
					Key:       "age",
					FieldType: NumberField,
					Required:  &TrueExpression,
				},
			},
			payload: map[string]interface{}{
				"age": 10,
			},
			assert: func(t *testing.T, form Form, submission FormSubmission, err error) {
				assert.NoError(t, err)
			},
		}, {
			name: "required numeric field: fail if payload property is absent",
			fields: []DataField{
				{
					Key:       "age",
					FieldType: NumberField,
					Required:  &TrueExpression,
				},
			},
			payload: map[string]interface{}{},
			assert: func(t *testing.T, form Form, submission FormSubmission, err error) {
				assert.Error(t, err)
			},
		}, {
			name: "required numeric field: fail if payload property is of wrong type",
			fields: []DataField{
				{
					Key:       "age",
					FieldType: NumberField,
					Required:  &TrueExpression,
				},
			},
			payload: map[string]interface{}{
				"age": "abc",
			},
			assert: func(t *testing.T, form Form, submission FormSubmission, err error) {
				assert.Error(t, err)
			},
		}, {
			name: "required numeric field with default value: should not trow if the payload property is missing",
			fields: []DataField{
				{
					Key:       "age",
					FieldType: NumberField,
					Required:  &TrueExpression,
					Default:   exprPtr("123"),
				},
			},
			payload: map[string]interface{}{},
			assert: func(t *testing.T, form Form, submission FormSubmission, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int64(123), submission.Payload["age"])
			},
		}, {
			name: "required single choice: should not throw if payload property is present",
			fields: []DataField{
				{
					Key:       "status",
					FieldType: ChoiceField,
					Options: []Choice{
						{
							Value: "A",
						}, {
							Value: "B",
						},
					},
					Required: &TrueExpression,
				},
			},
			payload: map[string]interface{}{
				"status": "A",
			},
			assert: func(t *testing.T, form Form, submission FormSubmission, err error) {
				assert.NoError(t, err)
			},
		}, {
			name: "required single choice: fail if the payload property is missing",
			fields: []DataField{
				{
					Key:       "status",
					FieldType: ChoiceField,
					Options: []Choice{
						{
							Value: "A",
						}, {
							Value: "B",
						},
					},
					Required: &TrueExpression,
				},
			},
			payload: map[string]interface{}{},
			assert: func(t *testing.T, form Form, submission FormSubmission, err error) {
				assert.Error(t, err)
			},
		}, {
			name: "required single choice: fail if the payload property is of wrong type",
			fields: []DataField{
				{
					Key:       "status",
					FieldType: ChoiceField,
					Options: []Choice{
						{
							Value: "A",
						}, {
							Value: "B",
						},
					},
					Required: &TrueExpression,
				},
			},
			payload: map[string]interface{}{
				"status": []int{1, 2, 4},
			},
			assert: func(t *testing.T, form Form, submission FormSubmission, err error) {
				assert.Error(t, err)
			},
		}, {
			name: "required single choice field with default: should not throw if the payload property is missing",
			fields: []DataField{
				{
					Key:       "status",
					FieldType: ChoiceField,
					Options: []Choice{
						{
							Value: "A",
						}, {
							Value: "B",
						},
					},
					Required: &TrueExpression,
					Default:  exprPtr("'A'"),
				},
			},
			payload: map[string]interface{}{},
			assert: func(t *testing.T, form Form, submission FormSubmission, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "A", submission.Payload["status"])
			},
		}, {
			name: "required reference field: should not throw if the payload property is present",
			fields: []DataField{
				{
					Key:       "ref",
					FieldType: ReferenceField,
					Required:  &TrueExpression,
					Default:   exprPtr("'A'"),
				},
			},
			payload: map[string]interface{}{
				"ref": map[string]interface{}{
					"kind":  "Bla",
					"group": "bla.com",
					"name":  "some id",
				},
			},
			assert: func(t *testing.T, form Form, submission FormSubmission, err error) {
				assert.NoError(t, err)
			},
		}, {
			name: "required reference field: fail if the payload property is absent",
			fields: []DataField{
				{
					Key:       "ref",
					FieldType: ReferenceField,
					Required:  &TrueExpression,
				},
			},
			payload: map[string]interface{}{},
			assert: func(t *testing.T, form Form, submission FormSubmission, err error) {
				assert.Error(t, err)
			},
		}, {
			name: "required reference field: fail if the payload reference field has invalid data",
			fields: []DataField{
				{
					Key:       "ref",
					FieldType: ReferenceField,
					Required:  &TrueExpression,
				},
			},
			payload: map[string]interface{}{
				"ref": map[string]interface{}{
					"a": "Bla",
					"b": "bla.com",
					"c": "some id",
				},
			},
			assert: func(t *testing.T, form Form, submission FormSubmission, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			f := Form{
				Fields: tc.fields,
			}
			submission := FormSubmission{
				Payload: tc.payload,
			}
			err := ValidateSubmission(submission, f)

			if tc.assert != nil {
				tc.assert(t, f, submission, err)
			}
		})
	}

}

func TestFieldAssignment(t *testing.T) {

	testCases := []struct {
		name      string
		key       string
		value     *Expression
		condition *Expression
		assert    func(t *testing.T, result map[string]interface{}, err error)
	}{
		{
			name:      "assign true boolean",
			key:       "abc",
			value:     exprPtr("true"),
			condition: nil,
			assert: func(t *testing.T, result map[string]interface{}, err error) {
				assert.NoError(t, err)
				assert.Equal(t, true, result["abc"])
			},
		}, {
			name:      "assign false boolean",
			key:       "abc",
			value:     exprPtr("false"),
			condition: nil,
			assert: func(t *testing.T, result map[string]interface{}, err error) {
				assert.NoError(t, err)
				assert.Equal(t, false, result["abc"])
			},
		}, {
			name:      "assign numeric",
			key:       "abc",
			value:     exprPtr("123"),
			condition: nil,
			assert: func(t *testing.T, result map[string]interface{}, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int64(123), result["abc"])
			},
		}, {
			name:      "assign string",
			key:       "abc",
			value:     exprPtr("'def'"),
			condition: nil,
			assert: func(t *testing.T, result map[string]interface{}, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "def", result["abc"])
			},
		}, {
			name:      "assign empty string",
			key:       "abc",
			value:     exprPtr("''"),
			condition: nil,
			assert: func(t *testing.T, result map[string]interface{}, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "", result["abc"])
			},
		}, {
			name:      "assign object json",
			key:       "abc",
			value:     exprPtr(`{"a":"b"}`),
			condition: nil,
			assert: func(t *testing.T, result map[string]interface{}, err error) {
				assert.NoError(t, err)
				assert.Equal(t, map[string]interface{}{
					"a": "b",
				}, result["abc"])
			},
		}, {
			name:      "assign list json",
			key:       "abc",
			value:     exprPtr(`[1, 2]`),
			condition: nil,
			assert: func(t *testing.T, result map[string]interface{}, err error) {
				assert.NoError(t, err)
				assert.Equal(t, []int64{1, 2}, result["abc"])
			},
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			val := map[string]interface{}{}
			assignment := FieldAssignment{Key: tc.key, Value: tc.value, Condition: tc.condition}
			err := assignment.Apply(val)
			tc.assert(t, val, err)
		})
	}
}

/**

Intake Form (i)
===============
FirstName
LastName
Age
LegalAssistance
ProtectionAssistance

Household (h)
=========
HeadOfHousehold
Members

Beneficiary (b)
===============
FirstName
LastName
Age

Service (s)
===========
Status (Pending, Opened, Closed)


1. User submits intake form
2. Beneficiary (b) is created

	b.FirstName = i.FirstName
	b.LastName = i.LastName
	b.Age = i.Age

3. Case is created

	case.Beneficiary = b
	case.Status = "Pending"

*/
