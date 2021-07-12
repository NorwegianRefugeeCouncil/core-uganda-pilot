package validator

type JSONValidationResponse struct {
	IsValid  bool       `json:"isValid"`
	Elements Validation `json:"elements"`
}
