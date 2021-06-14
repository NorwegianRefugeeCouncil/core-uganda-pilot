package v1

type Locale string

type Translation struct {
	Locale      Locale `json:"locale" bson:"locale" bson:"locale"`
	Translation string `json:"translation" bson:"translation" bson:"translation"`
}

type Translations []Translation
