package i18n

type LocaleString struct {
	Locale string `json:"locale" bson:"locale"`
	Value  string `json:"value" bson:"value"`
}

type Strings []LocaleString
