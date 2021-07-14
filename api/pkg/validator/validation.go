package validator

import (
	"encoding/json"
	"net/mail"
	"strconv"
	"unicode"
)

type Validate struct {
	Subject    map[string]interface{}
	Rules      ValidationRules
	IsValid    bool
	Validation Validation
}

type ValidationRules map[string]Rule

type ValueType int

const (
	Text         ValueType = iota // Text = 0
	Email        ValueType = iota // Email = 1
	Password     ValueType = iota // ...
	Name         ValueType = iota
	Phone        ValueType = iota
	UUID         ValueType = iota
	JSONTemplate ValueType = iota
)

type Rule struct {
	Type     ValueType
	Required bool
}

type Validation map[string]Status

type Status struct {
	ID      string      `json:"id"`
	Value   interface{} `json:"value"`
	IsValid bool        `json:"isValid"`
	Message string      `json:"message"`
}

func New(obj interface{}) *Validate {
	v := new(Validate)

	subject, err := structToMap(obj)
	if err != nil {
		panic(err)
	}

	v.Subject = subject
	v.IsValid = true
	return v
}

func structToMap(data interface{}) (map[string]interface{}, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	mapData := make(map[string]interface{})
	err = json.Unmarshal(dataBytes, &mapData)
	if err != nil {
		return nil, err
	}
	return mapData, nil
}

func (v *Validate) Validate(rules ValidationRules) bool {
	v.Rules = rules
	for field, rule := range v.Rules {
		// Check field
		value := v.Subject[field]
		msg := validateValue(value, rule)
		s := Status{
			ID:      field,
			Value:   value,
			IsValid: msg == "",
			Message: msg,
		}
		// Update
		v.Validation[field] = s
		if !s.IsValid {
			v.IsValid = false
		}
	}
	return v.IsValid
}

func (v *Validate) Response() ([]byte, error) {
	return json.Marshal(JSONValidationResponse{
		IsValid:  v.IsValid,
		Elements: v.Validation,
	})
}

func validateValue(value interface{}, rule Rule) string {
	str := value.(string)
	if rule.Required && str == "" {
		return "Field is required"
	}
	msg := ""
	switch rule.Type {
	case Text:
		if !IsAlphaNumeric(str) {
			msg = "Field contains invalid characters"
		}
	case Email:
		if !IsEmail(str) {
			msg = "Invalid email address"
		}
	case Password:
		if !IsPassword(str) {
			msg = "Invalid password"
		}
	case Name:
		if !IsAlpha(str) {
			msg = "Field contains invalid characters"
		}
	case Phone:
		if !IsPhone(str) {
			msg = "Invalid phone number"
		}
	case UUID:
		if !IsUUID(str) {
			msg = "Invalid UUID"
		}
	case JSONTemplate:
		if !IsJSON(str) {
			msg = "Invalid JSON template"
		}
	}
	return msg
}

func IsAlpha(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func IsNumeric(s string) bool {
	if _, err := strconv.Atoi(s); err != nil {
		return false
	}
	return true
}

func IsAlphaNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func IsJSON(s string) bool {
	return json.Valid([]byte(s))
}

func IsEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

func IsPassword(s string) bool {
	// TODO
	return true
}

func IsPhone(s string) bool {
	// TODO
	return true
}

func IsUUID(s string) bool {
	// TODO
	return true
}
