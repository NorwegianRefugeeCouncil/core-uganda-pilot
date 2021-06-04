package v1

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Service struct {
	ID        string
	ProjectID string
	Project   *Project
	Fields    []DataField
}

type FieldAssignment struct {
	Key       string
	Condition *Expression
	Value     *Expression
}

func (f FieldAssignment) Apply(into map[string]interface{}) error {
	if f.Condition != nil {
		condition, err := EvaluateBoolExpression(f.Condition)
		if err != nil {
			return err
		}
		if !condition {
			return nil
		}
	}
	value, ok, err := EvaluateValueExpression(f.Value)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	into[f.Key] = value
	return nil
}

type TriggerOperation struct {
	OperationType    string
	Kind             string
	Group            string
	Condition        *Expression
	FieldAssignments []FieldAssignment
}

type Trigger struct {
	Operations []TriggerOperation
}

type Beneficiary struct {
	ID          string
	HouseholdID *string
	Household   *Household
}

type Household struct {
	ID                string
	HeadOfHouseholdID *string
	HeadOfHousehold   Beneficiary
	Members           []*Beneficiary
}

type Case struct {
	ID            string
	BeneficiaryID string
	Beneficiary   *Beneficiary
	HouseholdID   string
	Household     *Household
	GroupID       string
	Group         *Group
}

type CaseFollowUp struct {
	ID     string
	CaseID string
	Case   *Case
}

type Group struct {
	ID string
}

type Membership struct {
	GroupID       string
	Group         *Group
	BeneficiaryID string
	Beneficiary   *Beneficiary
}

type Staff struct {
	ID string
}

type Project struct {
	ID string
}

type ExternalParty struct {
}

type SubmitForm struct {
}

type DataStructure struct {
}

type FieldType string

const (
	ShortText      FieldType = "shortText"
	LongText       FieldType = "longText"
	ChoiceField    FieldType = "choice"
	ReferenceField FieldType = "reference"
	ListField      FieldType = "list"
	NumberField    FieldType = "number"
	CheckBoxField  FieldType = "checkbox"
)

type Expression string

var TrueExpression Expression = "true"
var FalseExpression Expression = "false"

type Translation struct {
	Locale string
	Value  string
}

type TranslatedString []Translation

type Choice struct {
	Value string
	Label TranslatedString
}

type EntityType struct {
	Kind   string
	Group  string
	Filter *Expression
}

type DataField struct {
	FieldType  FieldType
	Key        string
	Immutable  bool
	Label      TranslatedString
	Tooltip    TranslatedString
	Help       TranslatedString
	MinLength  *Expression
	MaxLength  *Expression
	Required   *Expression
	Min        *Expression
	Max        *Expression
	MinItems   *Expression
	MaxItems   *Expression
	Default    *Expression
	Visible    *Expression
	ReadOnly   *Expression
	Disabled   *Expression
	Multiple   *Expression
	Unique     *Expression
	MultipleOf *Expression
	Pattern    *Expression
	Options    []Choice
	Children   []DataField
}

type Form struct {
	Fields []DataField
}

type FormSubmission struct {
	Payload map[string]interface{}
}

type Topology struct {
	Structures []Term
}

type Term struct {
	Name   string
	Parent string
}

var numberRegex = regexp.MustCompile("^[1-9][0-9]*$")

func isNumberChar(str string) bool {
	if numberRegex.MatchString(str) {
		_, _ = strconv.ParseInt(str, 10, 64)
		return true
	}
	return false
}

func EvaluateValueExpression(expression *Expression) (interface{}, bool, error) {
	if expression == nil {
		return nil, false, nil
	}
	exprStr := strings.TrimSpace(string(*expression))
	if len(exprStr) == 0 {
		return nil, false, nil
	}

	if exprStr == "true" {
		return true, true, nil
	}
	if exprStr == "false" {
		return false, true, nil
	}

	if exprStr[0:1] == "'" && exprStr[len(exprStr)-1:] == "'" {
		return exprStr[1 : len(exprStr)-1], true, nil
	}

	numberRegex := regexp.MustCompile("^[1-9][0-9]*$")
	if numberRegex.MatchString(exprStr) {
		num, _ := strconv.ParseInt(exprStr, 10, 64)
		return num, true, nil
	}

	if exprStr[0:1] == "{" {
		val := map[string]interface{}{}
		if err := json.Unmarshal([]byte(exprStr), &val); err != nil {
			return nil, false, err
		}
		return val, true, nil
	}

	if exprStr[0:1] == "[" {
		var val interface{}

		exprRunes := []rune(exprStr)
		for i := range exprRunes {
			charStr := exprRunes[i]
			if charStr == '"' {
				var val = []string{}
				if err := json.Unmarshal([]byte(exprStr), &val); err != nil {
					return nil, false, err
				}
				return val, true, nil
			}
			if charStr == '{' {
				var val []map[string]interface{}
				if err := json.Unmarshal([]byte(exprStr), &val); err != nil {
					return nil, false, err
				}
				return val, true, nil
			}
			if isNumberChar(string(charStr)) {
				var val []int64
				if err := json.Unmarshal([]byte(exprStr), &val); err != nil {
					return nil, false, err
				}
				return val, true, nil
			}
		}

		if err := json.Unmarshal([]byte(exprStr), &val); err != nil {
			return nil, false, err
		}
		return val, true, nil
	}

	return nil, false, nil
}

func EvaluateBoolExpression(expression *Expression) (bool, error) {
	if expression == nil {
		return false, nil
	}
	if *expression == "true" {
		return true, nil
	}
	if *expression == "false" {
		return false, nil
	}
	return false, fmt.Errorf("invalid boolean expression")
}

func isMultiple(submission FormSubmission, form Form, field DataField) (bool, error) {
	if field.Multiple == nil {
		return false, nil
	}
	isMultiple, err := EvaluateBoolExpression(field.Multiple)
	if err != nil {
		return false, err
	}
	return isMultiple, nil
}

func ValidateFieldRequired(submission FormSubmission, form Form, field DataField) error {

	isRequired, err := EvaluateBoolExpression(field.Required)
	if err != nil {
		return err
	}
	if !isRequired {
		return nil
	}
	fieldValue, hasFieldValue := submission.Payload[field.Key]
	if !hasFieldValue {
		return fmt.Errorf("field is required")
	}

	if fieldValue == nil {
		return fmt.Errorf("field is required")
	}

	switch t := fieldValue.(type) {
	case string:
		trimmed := strings.TrimSpace(t)
		if trimmed == "" {
			return fmt.Errorf("field is required")
		}
	}

	return nil
}

func ValidateFieldType(submission FormSubmission, form Form, field DataField) error {
	fieldValue, hasFieldValue := submission.Payload[field.Key]
	if !hasFieldValue {
		return nil
	}

	if fieldValue == nil {
		return nil
	}

	switch v := fieldValue.(type) {
	case string:

		// string values can be either
		// 1. shortText
		// 2. longText
		// 3. choice (no multiple)

		isMultiple, err := isMultiple(submission, form, field)
		if err != nil {
			return err
		}

		if field.FieldType != ShortText &&
			field.FieldType != LongText &&
			(field.FieldType != ChoiceField || isMultiple) {
			return fmt.Errorf("wrong field value")
		}
	case int, int8, int16, int32, int64, float32, float64:

		// a numeric value must be a numberField

		if field.FieldType != NumberField {
			return fmt.Errorf("wrong field value")
		}

	case []interface{}:

		// an array must be
		// 1. a list field
		// 2. a multiple choice field

		isMultiple, err := isMultiple(submission, form, field)
		if err != nil {
			return err
		}

		if field.FieldType != ListField && (field.FieldType != ChoiceField || !isMultiple) {
			return fmt.Errorf("invalid value type")
		}

	// TODO: validate item payload

	case map[string]interface{}:

		// an object can be a field of type
		// 1. object reference

		if field.FieldType != ReferenceField {
			return fmt.Errorf("invalid value type")
		}

		if field.FieldType == ReferenceField {
			kind, ok := v["kind"]
			if !ok {
				return fmt.Errorf("missing kind")
			}
			switch kindValue := kind.(type) {
			case string:
				if strings.TrimSpace(kindValue) == "" {
					return fmt.Errorf("kind cannot be empty")
				}
			default:
				return fmt.Errorf("invalid kind type")
			}

			group, ok := v["group"]
			if !ok {
				return fmt.Errorf("missing group")
			}
			switch groupValue := group.(type) {
			case string:
				if strings.TrimSpace(groupValue) == "" {
					return fmt.Errorf("kind cannot be empty")
				}
			default:
				return fmt.Errorf("invalid group type")
			}

			name, ok := v["name"]
			if !ok {
				return fmt.Errorf("missing name")
			}
			switch nameValue := name.(type) {
			case string:
				if strings.TrimSpace(nameValue) == "" {
					return fmt.Errorf("name cannot be empty")
				}
			default:
				return fmt.Errorf("invalid name type")
			}

		}

	default:
		return fmt.Errorf("unexpected field value")
	}
	return nil
}

func ValidateField(submission FormSubmission, form Form, field DataField) error {
	if err := ValidateFieldRequired(submission, form, field); err != nil {
		return err
	}
	if err := ValidateFieldType(submission, form, field); err != nil {
		return err
	}
	return nil
}

func PopulateFieldDefault(submission FormSubmission, form Form, field DataField) error {
	if field.Default == nil {
		return nil
	}
	_, hasFieldValue := submission.Payload[field.Key]
	if hasFieldValue {
		return nil
	}
	defaultValue, ok, err := EvaluateValueExpression(field.Default)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	submission.Payload[field.Key] = defaultValue
	return nil
}

func ValidateSubmission(submission FormSubmission, form Form) error {
	for _, field := range form.Fields {
		if err := PopulateFieldDefault(submission, form, field); err != nil {
			return err
		}
		if err := ValidateField(submission, form, field); err != nil {
			return err
		}
	}
	return nil
}

type UgandaIntakeForm struct {
	Form
}

type UgandaIntakeService struct {
}

type UgandaAssessmentService struct {
}

type UgandaAssessmentForm struct {
	Form
}

func NewUgandaIntakeForm() *UgandaAssessmentForm {
	return &UgandaAssessmentForm{
		Form: Form{
			Fields: []DataField{
				{
					Key:       "specificNeeds",
					FieldType: ChoiceField,
					Multiple:  &TrueExpression,
					Options: []Choice{
						{
							Value: "pregnantWoman",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Pregnant Woman",
								},
							},
						},
						{
							Value: "elderlyTakingCareOfMinorsAlone",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Elderly taking care of minors alone",
								},
							},
						},
						{
							Value: "singleParent",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Single parent",
								},
							},
						},
						{
							Value: "chronicIllness",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Chronic Illness",
								},
							},
						},
						{
							Value: "legalProtectionNeeds",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Legal Protection Needs",
								},
							},
						},
					},
				},
				{
					Key:       "disabilities",
					FieldType: ChoiceField,
					Multiple:  &TrueExpression,
					Options: []Choice{
						{
							Value: "moderatePhysicalImpairment",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Moderate Physical Impairment",
								},
							},
						},
						{
							Value: "severePhysicalImpairment",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Severe Physical Impairment",
								},
							},
						},
						{
							Value: "moderateSensoryImpairment",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Moderate Sensory Impairment",
								},
							},
						},
						{
							Value: "severeSensoryImpairment",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Severe Sensory Impairment",
								},
							},
						},
						{
							Value: "moderateMentalDisabilityOrIllness",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Moderate Mental Disability or Illness",
								},
							},
						},
						{
							Value: "severeMentalDisabilityOrIllness",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Severe Mental Disability or Illness",
								},
							},
						},
					},
				}, {
					Key:       "psychoSocialNeeds",
					FieldType: ChoiceField,
					Multiple:  &TrueExpression,
					Options: []Choice{
						{
							Value: "noNeeds",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "No needs",
								},
							},
						}, {
							Value: "substanceAbuse",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Substance Abuse",
								},
							},
						}, {
							Value: "sadnessOrDepression",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Sadness or Depression",
								},
							},
						}, {
							Value: "anxietyOrNervousness",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Anxiety or Nervousness",
								},
							},
						}, {
							Value: "angerManagementIssues",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Anger Management Issues",
								},
							},
						}, {
							Value: "sleepIssues",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Sleep Issues (lack of/excess of)",
								},
							},
						},
					},
				}, {
					Key:       "needsAccessAssessment",
					FieldType: ChoiceField,

					Options: []Choice{
						{
							Value: "canMeetNeedsWithoutWorry",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "We can meet our needs without worry",
								},
							},
						}, {
							Value: "canMeetNeeds",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "We can meet our needs",
								},
							},
						}, {
							Value: "barelyCanMeetNeeds",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "We can barely meet our needs",
								},
							},
						}, {
							Value: "unableToMeetNeeds",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "We are unable to meet our needs",
								},
							},
						}, {
							Value: "totallyUnableToMeetNeeds",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "We are totally unable to meet our needs",
								},
							},
						},
					},
				}, {
					Key:       "accessObstacles",
					FieldType: ChoiceField,
					Multiple:  &TrueExpression,
					Options: []Choice{
						{
							Value: "insufficientFunds",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Insufficient funds to buy / access goods or services",
								},
							},
						}, {
							Value: "travelDistance",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Travel distance makes it difficult to access markets / service providers",
								},
							},
						}, {
							Value: "insecurity",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Insecurity makes access to markets / service providers difficult",
								},
							},
						}, {
							Value: "socialDiscrimination",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Social discrimination makes it difficult to access markets / service providers",
								},
							},
						}, {
							Value: "insufficientQuantityOrQuality",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Insufficient quantity of goods / services available",
								},
							},
						}, {
							Value: "inadequateQuality",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Inadequate quality of goods / services",
								},
							},
						}, {
							Value: "insufficientCapacitiesOrIncompetence",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Insufficient capacities and incompetence of service providers",
								},
							},
						},
					},
				}, {
					Key:       "recommendations",
					FieldType: ChoiceField,
					Multiple:  &TrueExpression,
					Options: []Choice{
						{
							Value: "referral",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Referral",
								},
							},
						}, {
							Value: "inKindAssistance",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "In-Kind assistance",
								},
							},
						}, {
							Value: "cashAssistance",
							Label: []Translation{
								{
									Locale: "en",
									Value:  "Cash Assistance",
								},
							},
						},
					},
				},
			},
		},
	}
}
