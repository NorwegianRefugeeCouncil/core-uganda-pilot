package i18n

type LocaleString struct {
	Locale string `json:"locale" bson:"locale"`
	Value  string `json:"value" bson:"value"`
}

type Strings []LocaleString

func (s Strings) ForLocale(locale string) string {
	for _, localeString := range s {
		if localeString.Locale == locale {
			return localeString.Value
		}
	}
	return ""
}
