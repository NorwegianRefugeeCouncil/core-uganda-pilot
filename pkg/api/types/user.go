package types

type User struct {
	ID string `json:"id"`
}

type UserProfile struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
}
